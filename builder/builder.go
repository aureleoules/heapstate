package builder

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/shared"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/mholt/archiver"
	"github.com/phayes/freeport"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CloneRepository util
func CloneRepository(url string, branch string, dir string) (*git.Repository, error) {
	log.Println("BRANCH=", branch)
	fmt.Println("CLONING", url, "in", dir)
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
	})

	return r, err
}

func MakeTarball(dir string, dest string) error {
	z := archiver.Tar{
		MkdirAll:               true,
		ContinueOnError:        false,
		OverwriteExisting:      true,
		ImplicitTopLevelFolder: true,
	}

	var filesPath []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		p := dir + "/" + f.Name()
		filesPath = append(filesPath, p)
	}

	return z.Archive(filesPath, dest)
}

// Clean repo dir
func Clean() {
	dir := os.Getenv("CLONE_REPO_DIR")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		os.RemoveAll(dir + "/" + f.Name())
	}
}

// Build app
func Build(app shared.App) error {
	build := shared.Build{
		AppID:     app.ID,
		CreatedAt: time.Now(),
		Branch:    app.BuildOptions.Branch,
		Logs:      []string{},
	}

	build.Status = shared.Building
	build.Create()
	app.SetState(shared.Stopped)

	Clean()
	/* Clone repo */
	repoDir := os.Getenv("CLONE_REPO_DIR") + "/" + app.Name
	r, err := CloneRepository(app.CompleteURL, app.BuildOptions.Branch, repoDir)
	if err != nil {
		build.SetStatus(shared.BuildError, "Could not clone repository.")

		return err
	}

	ref, err := r.Head()
	if err != nil {
		build.SetStatus(shared.BuildError, "Could not get HEAD.")

		return err
	}

	commit, err := r.CommitObject(ref.Hash())

	if err != nil {
		build.SetStatus(shared.BuildError, "Could not get latest commit.")
		return err
	}

	build.SetCommit(commit.Hash.String(), commit.Message)

	/* Build tarball */
	tarballDir := os.Getenv("TARBALL_DIR") + "/" + app.Name + ".tar"
	err = MakeTarball(repoDir, tarballDir)
	if err != nil {
		log.Println(err)
		build.SetStatus(shared.BuildError, "Could not build tarball.")

		return err
	}

	/* Open tarball */
	dockerBuildContext, err := os.Open(tarballDir)
	if err != nil {

		build.SetStatus(shared.BuildError, "Could not open tarball.")

		return err
	}
	defer dockerBuildContext.Close()

	opt := types.ImageBuildOptions{
		SuppressOutput: false,
		PullParent:     true,
		Dockerfile:     "Dockerfile",
		Tags:           []string{app.Name},
	}
	resp, err := common.DockerClient.ImageBuild(context.Background(), dockerBuildContext, opt)
	if err != nil {
		fmt.Println(err)
		build.SetStatus(shared.BuildError, "Could not build Docker image.")

		return err
	}

	reg := regexp.MustCompile(`":\"(.*)\"}`)
	defer resp.Body.Close()
	rd := bufio.NewReader(resp.Body)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			build.SetStatus(shared.BuildError, "Unexpected error")
			return err
		}

		matches := reg.FindStringSubmatch(string(n))
		if len(matches) < 2 {
			continue
		}

		str := utils.UnescapeString(matches[1])
		if str == "\n" {
			continue
		}

		build.Log(str)
	}

	build.SetStatus(shared.Building, "")

	/* Delete previous docker container */
	err = common.DockerClient.ContainerRemove(context.Background(), app.ContainerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not kill container,", err)
	}

	exposedPorts := map[nat.Port]struct{}{"80/tcp": {}}

	port, err := freeport.GetFreePort()
	if err != nil {
		build.SetStatus(shared.BuildError, "Could not find available port.")
		return err
	}

	portStr := strconv.Itoa(port)

	netConfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"traefik_default": &network.EndpointSettings{
				NetworkID: os.Getenv("TRAEFIK_NETWORK_ID"),
			},
		},
	}
	dockerResponse, err := common.DockerClient.ContainerCreate(context.Background(), &container.Config{
		Image:        app.Name,
		ExposedPorts: exposedPorts,
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers." + app.Name + ".rule":             "Host(`" + app.Name + ".heapstate.com`)",
			"traefik.http.routers." + app.Name + ".entrypoints":      "websecure",
			"traefik.http.routers." + app.Name + ".tls.certresolver": "myresolver",
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port("80/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portStr}},
		},
		Resources: container.Resources{
			Memory: app.ContainerOptions.MaxRAM,
		},
	}, &netConfig, app.Name)

	if err != nil {
		build.SetStatus(shared.DeployError, "Could not deploy Docker container.")
		fmt.Println(err)
		return err
	}

	err = app.SetContainerID(dockerResponse.ID)
	if err != nil {
		build.SetStatus(shared.BuildError, "Could not set container ID.")

		return err
	}
	log.Println("Running", dockerResponse.ID)

	err = common.DockerClient.ContainerStart(context.Background(), dockerResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		build.SetStatus(shared.DeployError, "Could not start Docker container.")

		return err
	}

	build.SetStatus(shared.Deployed, "")
	app.SetState(shared.Running)

	return nil
}
