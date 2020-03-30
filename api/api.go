package api

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/heapstate/api/apps"
	"github.com/aureleoules/heapstate/api/users"
	usersHandlers "github.com/aureleoules/heapstate/handlers/users"

	"github.com/aureleoules/heapstate/api/oauth"
	"github.com/aureleoules/heapstate/models"
	"github.com/aureleoules/heapstate/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func handleProtected(r *gin.RouterGroup) {
	users.HandleProtected(r.Group("/users"))
	oauth.HandleProtected(r.Group("/oauth"))
	apps.HandleProtected(r.Group("/apps"))
}

func handlePub(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		utils.Response(c, http.StatusOK, nil, "Welcome to heapstate's API.")
	})

	users.HandlePub(r.Group("/users"))
}

var publicApi *gin.RouterGroup
var api *gin.RouterGroup

var authMiddleware *jwt.GinJWTMiddleware

var version = "v1"

func creater() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	publicApi = r.Group("/")
	api = r.Group("/")

	authMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "heapstate",
		Key:        []byte(os.Getenv("SECRET")),
		Timeout:    time.Hour * 12,
		MaxRefresh: time.Hour * 12,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":       v.ID,
					"username": v.Username,
					"email":    v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id, _ := primitive.ObjectIDFromHex(claims["id"].(string))
			return &models.User{
				ID: id,
			}
		},
		Authenticator: usersHandlers.Authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			_, err := primitive.ObjectIDFromHex(claims["id"].(string))
			if err != nil {
				return false
			}
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	api.Use(authMiddleware.MiddlewareFunc())
	{
		handlePub(publicApi)
		handleProtected(api)
	}
	publicApi.POST("/authenticate", authMiddleware.LoginHandler)

	return r
}

// Listen creates web r
func Listen(port string) {
	r := creater()
	r.Run(":" + port)
}
