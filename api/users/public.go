package users

import (
	"github.com/aureleoules/heapstate/handlers/users"
	"github.com/gin-gonic/gin"
)

// HandlePub hanldes public routes
func HandlePub(r *gin.RouterGroup) {
	r.POST("/", users.HandleRegister)
}
