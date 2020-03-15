package apps

import (
	"net/http"

	"github.com/aureleoules/heapstack/builder"
	"github.com/aureleoules/heapstack/shared"
	"github.com/aureleoules/heapstack/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func newAppHandler(c *gin.Context) {
	var app shared.App
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
	app.CompleteURL = "https://" + app.URL

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

func fetchAppHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app)
	return
}

func fetchBuildOptionsHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	utils.Response(c, http.StatusOK, nil, app.BuildOptions)
	return
}

func deployHandler(c *gin.Context) {
	name := c.Param("name")

	app, err := FetchApp(name)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}

	spew.Dump(app)

	builder.Build(app)
}

func fetchBuildsHandler(c *gin.Context) {
	name := c.Param("name")

	id, err := GetAppID(name)
	if err != nil {
		utils.Response(c, http.StatusNotFound, err, nil)
		return
	}

	builds, err := GetBuilds(id)
	if err != nil {
		utils.Response(c, http.StatusInternalServerError, err, nil)
		return
	}
	utils.Response(c, http.StatusOK, nil, builds)
	return
}
