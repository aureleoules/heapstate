package apps

import "github.com/gin-gonic/gin"

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/", newAppHandler)
	r.GET("/:name", fetchAppHandler)
	r.GET("/", fetchAppsHandler)

	r.GET("/:name/buildoptions", fetchBuildOptionsHandler)
	r.POST("/:name/deploy", deployHandler)

	r.GET("/:name/builds", fetchBuildsHandler)
	r.GET("/:name/builds/:id", fetchBuildHandler)
}
