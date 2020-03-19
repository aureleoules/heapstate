package apps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aureleoules/heapstate/builder"
	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/shared"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newAppHandler(c *gin.Context) {
	var app shared.App
	err := c.BindJSON(&app)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	err = app.Validate()
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, err, nil)
		return
	}

	var baseURL string
	switch app.Provider {
	case shared.GitHubProvider:
		baseURL = "github.com/"
		break
	case shared.GitLabProvider:
		baseURL = "gitlab.com/"
		break
	case shared.BitBucketProvider:
		baseURL = "bitbucket.org/"
	}

	// Set app repo url
	app.URL = baseURL + app.Owner + "/" + app.Name
	app.CompleteURL = "https://" + app.URL

	app.UserID = utils.ExtractUserID(c)

	err = app.Save()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, nil)
}

func fetchAppsHandler(c *gin.Context) {
	userID := utils.ExtractUserID(c)

	apps, err := FetchApps(userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	for i := range apps {
		apps[i].LastBuild, err = GetLatestBuild(apps[i].ID)
	}

	utils.Response(c, http.StatusOK, nil, apps)
	return
}

func fetchStats(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	response, err := common.DockerClient.ContainerStats(context.Background(), app.ContainerID, false)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	var stats types.StatsJSON
	err = json.Unmarshal(buf.Bytes(), &stats)

	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, shared.ContainerStats{
		RAMUsage: float64(stats.MemoryStats.Usage+stats.MemoryStats.Stats["cache"]) / 1024 / 1024 / 8,
		MaxRAM:   app.ContainerOptions.MaxRAM,
		CPU:      utils.CalculateCPUPercentUnix(0, 0, &stats),
	})

}

func fetchAppHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app)
	return
}

func fetchBuildOptionsHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app.BuildOptions)
	return
}

func deployHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	builder.Build(app)
}

func fetchBuildsHandler(c *gin.Context) {
	name := c.Param("name")

	limitStr := c.DefaultQuery("limit", "4")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, errors.New("invalid limit"), nil)
		return
	}

	id, err := GetAppID(name)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
		return
	}

	builds, err := GetBuilds(id, limit)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	utils.Response(c, http.StatusOK, nil, builds)
	return
}

func fetchBuildHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, errors.New("invalid id"), nil)
		return
	}

	build, err := GetBuild(objectID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	utils.Response(c, http.StatusOK, nil, build)
	return
}

func fetchLogsHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	response, err := common.DockerClient.ContainerLogs(context.Background(), app.ContainerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response)

	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, buf.String())
}

func fetchContainerOptionsHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app.ContainerOptions)
	return
}
