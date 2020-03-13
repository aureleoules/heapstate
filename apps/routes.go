package apps

import "github.com/gin-gonic/gin"

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/", newAppHandler)
	r.GET("/", fetchAppsHandler)
}
