package router

import (
	"github.com/aureleoules/heapstack/users"
	"github.com/gin-gonic/gin"
)

func handlePub(r *gin.RouterGroup) {
	users.HandlePub(r.Group("/users"))
}
