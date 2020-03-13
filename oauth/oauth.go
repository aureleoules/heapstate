package oauth

import "github.com/gin-gonic/gin"

// HandleProtected oauth routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/github", exchangeGitHubTokenHandler)
}
