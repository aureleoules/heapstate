package users

import (
	"github.com/aureleoules/heapstate/handlers/users"
	"github.com/gin-gonic/gin"
)

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.GET("/profile", users.HandleGetProfile)
}
