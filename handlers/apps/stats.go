package apps

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/models"
	"github.com/aureleoules/heapstate/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

func FetchStatsHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
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

	utils.Response(c, http.StatusOK, nil, models.ContainerStats{
		RAMUsage: int64(stats.MemoryStats.Usage + stats.MemoryStats.Stats["cache"]),
		MaxRAM:   app.ContainerOptions.MaxRAM,
		CPU:      utils.CalculateCPUPercentUnix(0, 0, &stats),
	})

}
