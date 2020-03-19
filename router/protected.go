package router

import (
	"github.com/aureleoules/heapstate/apps"
	"github.com/aureleoules/heapstate/oauth"
	"github.com/aureleoules/heapstate/users"
	"github.com/gin-gonic/gin"
)

func handleProtected(r *gin.RouterGroup) {
	users.HandleProtected(r.Group("/users"))
	oauth.HandleProtected(r.Group("/oauth"))
	apps.HandleProtected(r.Group("/apps"))

}
