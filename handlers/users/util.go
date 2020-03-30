package users

import (
	"context"

	"github.com/aureleoules/heapstate/common"
	"github.com/aureleoules/heapstate/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// RetrieveByEmail finds user by email
func RetrieveByEmail(email string) (*models.User, error) {
	result := common.DB.Collection(common.UsersCollection).FindOne(context.Background(), bson.M{
		"email": email,
	})

	var user models.User
	err := result.Decode(&user)

	return &user, err
}

// GetUser by ID
func GetUser(id primitive.ObjectID) (*models.User, error) {
	r := common.DB.Collection(common.UsersCollection).FindOne(context.Background(), bson.M{
		"_id": id,
	})
	var user models.User
	err := r.Decode(&user)
	return &user, err
}
