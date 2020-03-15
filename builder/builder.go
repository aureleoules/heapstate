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
	"time"

	"github.com/aureleoules/heapstack/shared"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/mholt/archiver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("Initialized Docker client")

}

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
		log.Println(f.Name())
		err := os.RemoveAll(dir + "/" + f.Name())
		log.Println(err)
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
	resp, err := cli.ImageBuild(context.Background(), dockerBuildContext, opt)
	if err != nil {
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
		log.Println(string(n))
		log.Println(matches)
		if len(matches) < 2 {
			continue
		}
		log.Println(matches[1])
		build.Log(matches[1])
	}

	build.SetStatus(shared.Building, "")

	/* Delete previous docker container */
	err = cli.ContainerRemove(context.Background(), app.ContainerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Println("Could not kill container,", err)
	}

	exposedPorts := map[nat.Port]struct{}{"80/tcp": {}}
	dockerResponse, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        app.Name,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port("80/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5000"}},
		},
	}, &network.NetworkingConfig{}, app.Name)

	err = app.SetContainerID(dockerResponse.ID)
	if err != nil {
		build.SetStatus(shared.BuildError, "Could not set container ID.")

		return err
	}

	if err != nil {
		build.SetStatus(shared.DeployError, "Could not deploy Docker container.")

		return err
	}
	log.Println("Running", dockerResponse.ID)

	err = cli.ContainerStart(context.Background(), dockerResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		build.SetStatus(shared.DeployError, "Could not start Docker container.")

		return err
	}

	build.SetStatus(shared.Deployed, "")

	return nil
}
