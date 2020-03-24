package router

import (
	"net/http"

	"github.com/aureleoules/heapstate/users"
	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-gonic/gin"
)

func handlePub(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		utils.Response(c, http.StatusOK, nil, "Welcome to heapstate's API.")
	})

	users.HandlePub(r.Group("/users"))

}
