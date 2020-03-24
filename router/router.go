package router

import (
	"os"
	"time"

	"github.com/aureleoules/heapstate/users"
	"github.com/gin-contrib/cors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
			if v, ok := data.(*users.User); ok {
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
			return &users.User{
				ID: id,
			}
		},
		Authenticator: users.Authenticator,
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
