package utils

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExtractUserID from JWT token
func ExtractUserID(c *gin.Context) primitive.ObjectID {
	claims := jwt.ExtractClaims(c)
	id, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	return id
}
