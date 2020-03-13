package users

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// HandlePub hanldes public routes
func HandlePub(r *gin.RouterGroup) {
	r.POST("/", handleRegister)
}

// HandleProtected routes
func HandleProtected(r *gin.RouterGroup) {

}

// Authenticator handler
func Authenticator(c *gin.Context) (interface{}, error) {
	log.Println("Authenticating")
	var user User
	c.BindJSON(&user)
	u := RetrieveByEmail(user.Email)
	if u == nil {
		return nil, jwt.ErrFailedAuthentication
	}

	log.Print("ok")

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
