package apps

import (
	"net/http"

	"github.com/aureleoules/heapstate/models"
	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-gonic/gin"
)

func FetchBuildOptionsHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app.BuildOptions)
	return
}

func SaveBuildOptionsHandler(c *gin.Context) {
	name := c.Param("name")

	var options models.BuildOptions
	err := c.BindJSON(&options)
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, err, nil)
		return
	}
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
		return
	}
	err = app.SaveBuildOptions(options)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, nil)
	return
}
