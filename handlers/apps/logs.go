package apps

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

func FetchLogsHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
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
