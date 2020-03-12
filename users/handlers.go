package users

import (
	"net/http"

	"github.com/aureleoules/heapstack/utils"
	"github.com/gin-gonic/gin"
)

func handleAuthenticate(c *gin.Context) {

}

func handleRegister(c *gin.Context) {
	var user User
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
