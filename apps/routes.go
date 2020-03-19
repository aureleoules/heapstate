package apps

import "github.com/gin-gonic/gin"

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/", newAppHandler)
	r.GET("/", fetchAppsHandler)

	r.GET("/:name", fetchAppHandler)
	r.GET("/:name/stats", fetchStats)
	r.GET("/:name/buildoptions", fetchBuildOptionsHandler)
	r.GET("/:name/builds", fetchBuildsHandler)
	r.GET("/:name/builds/:id", fetchBuildHandler)
	r.GET("/:name/logs", fetchLogsHandler)
	r.GET("/:name/containeroptions", fetchContainerOptionsHandler)
	r.POST("/:name/deploy", deployHandler)

}
