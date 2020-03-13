package oauth

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aureleoules/heapstack/utils"
	"github.com/gin-gonic/gin"
)

const (
	GITHUB_ENDPOINT = "https://github.com/login/oauth/access_token"
)

type Request struct {
	Code string `json:"code"`
}

func ExchangeGitHubTokenHandler(c *gin.Context) {
	var req Request
	c.BindJSON(&req)

	resp, err := http.PostForm(GITHUB_ENDPOINT, url.Values{
		"client_id":     {os.Getenv("GITHUB_ID")},
		"client_secret": {os.Getenv("GITHUB_SECRET")},
		"code":          {req.Code},
	})
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	m, _ := url.ParseQuery(string(body))

	log.Print(string(body))
	token := m.Get("access_token")
	log.Println(token)

	if err != nil {
		log.Print(err)
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, token)
	return
}
