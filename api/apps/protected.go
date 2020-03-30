package apps

import (
	apps "github.com/aureleoules/heapstate/handlers/apps"
	"github.com/gin-gonic/gin"
)

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/", apps.NewAppHandler)
	r.GET("/", apps.FetchAppsHandler)

	r.GET("/:name", apps.FetchAppHandler)
	r.GET("/:name/stats", apps.FetchStatsHandler)

	r.POST("/:name/start", apps.StartHandler)
	r.POST("/:name/restart", apps.RestartHandler)
	r.POST("/:name/stop", apps.StopHandler)

	r.GET("/:name/buildoptions", apps.FetchBuildOptionsHandler)
	r.PUT("/:name/buildoptions", apps.SaveBuildOptionsHandler)

	r.GET("/:name/builds", apps.FetchBuildsHandler)
	r.GET("/:name/builds/:id", apps.FetchBuildHandler)

	r.GET("/:name/logs", apps.FetchLogsHandler)

	r.GET("/:name/containeroptions", apps.FetchContainerOptionsHandler)
	r.PUT("/:name/containeroptions", apps.SaveContainerOptionsHandler)

	r.POST("/:name/deploy", apps.DeployHandler)
}
