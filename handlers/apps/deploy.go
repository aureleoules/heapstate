package apps

import (
	"net/http"

	"github.com/aureleoules/heapstate/builder"
	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-gonic/gin"
)

func DeployHandler(c *gin.Context) {
	name := c.Param("name")
	userID := utils.ExtractUserID(c)

	app, err := FetchApp(name, userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	builder.Build(app)
}
