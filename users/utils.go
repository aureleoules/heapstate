package users

import (
	"context"

	"github.com/aureleoules/heapstack/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// RetrieveByEmail finds user by email
func RetrieveByEmail(email string) (*User, error) {
	result := common.DB.Collection(common.UsersCollection).FindOne(context.Background(), bson.M{
		"email": email,
	})

	var user User
	err := result.Decode(&user)

	return &user, err
}

// GetUser by ID
func GetUser(id primitive.ObjectID) (*User, error) {
	r := common.DB.Collection(common.UsersCollection).FindOne(context.Background(), bson.M{
		"_id": id,
	})
	var user User
	err := r.Decode(&user)
	return &user, err
}
