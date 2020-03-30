package oauth

import (
	"github.com/aureleoules/heapstate/handlers/oauth"
	"github.com/gin-gonic/gin"
)

// HandleProtected oauth routes
func HandleProtected(r *gin.RouterGroup) {
	r.POST("/github", oauth.ExchangeGitHubTokenHandler)
}
