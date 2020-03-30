package users

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/aureleoules/heapstate/models"
	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleGetProfile(c *gin.Context) {
	id := utils.ExtractUserID(c)
	log.Println("extract", id)

	user, err := GetUser(id)
	if err != nil {
		utils.Response(c, http.StatusUnauthorized, nil, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, user.Public())
	return
}

func HandleRegister(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)

	err := user.Validate()
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, err, nil)
		return
	}

	user.HashPassword()

	err = user.Save()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, user.Public())
	return
}

// Authenticator handler
func Authenticator(c *gin.Context) (interface{}, error) {
	var user models.User
	c.BindJSON(&user)
	u, err := RetrieveByEmail(user.Email)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &models.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}
