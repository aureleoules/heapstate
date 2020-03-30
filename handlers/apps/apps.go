package apps

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aureleoules/heapstate/builder"
	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/models"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/taion809/haikunator"
)

func NewAppHandler(c *gin.Context) {
	var app models.App
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

func FetchAppsHandler(c *gin.Context) {
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

func FetchAppHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app)
	return
}

func StartHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
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

func RestartHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	app.SetState(models.Stopped)

	timeout := time.Duration(10 * time.Second)
	err = common.DockerClient.ContainerRestart(context.Background(), app.ContainerID, &timeout)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	app.SetState(models.Running)

	utils.Response(c, http.StatusOK, nil, nil)
	return
}

func StopHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
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

	app.SetState(models.Stopped)
	utils.Response(c, http.StatusOK, nil, nil)

	return
}
