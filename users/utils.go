package users

import (
	"context"

	"github.com/aureleoules/heapstack/common"
	"gopkg.in/mgo.v2/bson"
)

// RetrieveByEmail finds user by email
func RetrieveByEmail(email string) *User {
	result := common.DB.Collection(common.UsersCollection).FindOne(context.Background(), bson.M{
		"email": email,
	})

	var user User
	result.Decode(&user)

	return &user
}
