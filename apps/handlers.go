package apps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aureleoules/heapstate/builder"
	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/shared"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/taion809/haikunator"

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

	app.URL = strings.ReplaceAll(app.CompleteURL, "https://", "")

	app.UserID = utils.ExtractUserID(c)

	h := haikunator.NewHaikunator()
	app.Name = h.TokenHaikunate(10000)

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

func fetchStatsHandler(c *gin.Context) {
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
		RAMUsage: int64(stats.MemoryStats.Usage + stats.MemoryStats.Stats["cache"]),
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

func startHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	err = builder.Build(app)

	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, nil)
	return
}

func restartHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	app.SetState(shared.Stopped)

	timeout := time.Duration(10 * time.Second)
	err = common.DockerClient.ContainerRestart(context.Background(), app.ContainerID, &timeout)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	app.SetState(shared.Running)

	utils.Response(c, http.StatusOK, nil, nil)
	return
}

func stopHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	err = common.DockerClient.ContainerRemove(context.Background(), app.ContainerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	app.SetState(shared.Stopped)
	utils.Response(c, http.StatusOK, nil, nil)

	return
}

func saveContainerOptionsHandler(c *gin.Context) {
	name := c.Param("name")

	var options shared.ContainerOptions
	err := c.BindJSON(&options)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	err = app.SaveContainerOptions(options)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, nil)
	return
}
