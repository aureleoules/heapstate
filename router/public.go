package router

import (
	"github.com/aureleoules/heapstate/users"
	"github.com/gin-gonic/gin"
)

func handlePub(r *gin.RouterGroup) {
	users.HandlePub(r.Group("/users"))
}
