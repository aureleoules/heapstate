package users

import (
	"github.com/gin-gonic/gin"
)

// HandlePub hanldes public routes
func HandlePub(r *gin.RouterGroup) {
	r.POST("/", handleRegister)

}

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {
	r.GET("/profile", handleGetProfile)
}
