package apps

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchBuildsHandler(c *gin.Context) {
	name := c.Param("name")

	limitStr := c.DefaultQuery("limit", "4")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, errors.New("invalid limit"), nil)
		return
	}
	userID := utils.ExtractUserID(c)

	id, err := GetAppID(name, userID)
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

func FetchBuildHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, errors.New("invalid id"), nil)
		return
	}
	userID := utils.ExtractUserID(c)

	build, err := GetBuild(objectID, userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	utils.Response(c, http.StatusOK, nil, build)
	return
}
