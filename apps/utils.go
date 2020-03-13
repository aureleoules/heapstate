package apps

import (
	"context"

	"github.com/aureleoules/heapstack/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"labix.org/v2/mgo/bson"
)

// FetchApps of user
func FetchApps(userID primitive.ObjectID) ([]App, error) {
	c, err := common.DB.Collection(common.AppsCollection).Find(context.Background(), bson.M{
		"user_id": userID,
	})

	if err != nil {
		return nil, err
	}
	var apps []App

	err = c.All(context.Background(), &apps)
	return apps, nil
}
