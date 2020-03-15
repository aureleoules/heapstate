package apps

import (
	"context"

	"github.com/aureleoules/heapstack/common"
	"github.com/aureleoules/heapstack/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"labix.org/v2/mgo/bson"
)

// FetchApps of user
func FetchApps(userID primitive.ObjectID) ([]shared.App, error) {
	c, err := common.DB.Collection(common.AppsCollection).Find(context.Background(), bson.M{
		"user_id": userID,
	})

	if err != nil {
		return nil, err
	}
	var apps []shared.App

	err = c.All(context.Background(), &apps)
	return apps, nil
}

// FetchApp returns single app by name
func FetchApp(userID primitive.ObjectID, name string) (shared.App, error) {
	r := common.DB.Collection(common.AppsCollection).FindOne(context.Background(), bson.M{
		"user_id": userID,
		"name":    name,
	})

	var app shared.App
	err := r.Decode(&app)
	if err != nil {
		return shared.App{}, err
	}

	build, err := GetLatestBuild(app.ID)
	app.LastBuild = build

	return app, err
}

// GetBuilds of app
func GetBuilds(appID primitive.ObjectID) ([]shared.Build, error) {
	r, err := common.DB.Collection(common.BuildsCollection).Find(context.Background(), bson.M{
		"app_id": appID,
	})
	if err != nil {
		return nil, err
	}
	var builds []shared.Build
	err = r.Decode(&builds)
	return builds, err
}

// GetLatestBuild of app
func GetLatestBuild(appID primitive.ObjectID) (shared.Build, error) {
	r := common.DB.Collection(common.BuildsCollection).FindOne(context.Background(), bson.M{
		"app_id": appID,
	})
	var build shared.Build
	err := r.Decode(&build)
	return build, err
}