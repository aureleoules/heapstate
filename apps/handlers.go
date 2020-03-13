package apps

import (
	"net/http"

	"github.com/aureleoules/heapstack/shared"
	"github.com/aureleoules/heapstack/utils"
	"github.com/gin-gonic/gin"
)

func newAppHandler(c *gin.Context) {
	var app App
	err := c.BindJSON(&app)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	err = app.Validate()
	if err != nil {
		utils.Response(c, http.StatusNotAcceptable, err, nil)
		return
	}

	var baseURL string
	switch app.Provider {
	case shared.GitHubProvider:
		baseURL = "github.com/"
		break
	case shared.GitLabProvider:
		baseURL = "gitlab.com/"
		break
	case shared.BitBucketProvider:
		baseURL = "bitbucket.org/"
	}

	// Set app repo url
	app.URL = baseURL + app.Owner + "/" + app.Name

	app.UserID = utils.ExtractUserID(c)

	err = app.Save()
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, nil)
}

func fetchAppsHandler(c *gin.Context) {
	userID := utils.ExtractUserID(c)

	apps, err := FetchApps(userID)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, apps)
	return
}
