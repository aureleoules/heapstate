package router

import (
	"github.com/aureleoules/heapstack/oauth"
	"github.com/aureleoules/heapstack/users"
	"github.com/gin-gonic/gin"
)

func handleProtected(r *gin.RouterGroup) {
	users.HandleProtected(r.Group("/users"))
	oauth.HandleProtected(r.Group("/oauth"))

}
