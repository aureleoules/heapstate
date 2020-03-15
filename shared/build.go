package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/aureleoules/heapstack/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Build struct
type Build struct {
	ID    primitive.ObjectID `json:"-" bson:"_id"`
	AppID primitive.ObjectID `json:"app_id" bson:"app_id"`

	Branch        string   `json:"branch" bson:"branch"`
	CommitHash    string   `json:"commit_hash" bson:"commit_hash"`
	CommitMessage string   `json:"commit_message" bson:"commit_message"`
	Logs          []string `json:"logs" bson:"logs"`

	Status DeployStatus `json:"status" bson:"status"`
	Error  string       `json:"error" bson:"error"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Create build
func (b *Build) Create() error {
	b.ID = primitive.NewObjectID()
	_, err := common.DB.Collection(common.BuildsCollection).InsertOne(context.Background(), b)

	return err
}

// SetCommit of build
func (b *Build) SetCommit(hash string, message string) error {
	_, err := common.DB.Collection(common.BuildsCollection).UpdateOne(context.Background(), bson.M{
		"_id": b.ID,
	}, bson.M{
		"$set": bson.M{
			"commit_hash":    hash,
			"commit_message": message,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// SetStatus of build
func (b *Build) SetStatus(status DeployStatus, errorMessage string) error {
	_, err := common.DB.Collection(common.BuildsCollection).UpdateOne(context.Background(), bson.M{
		"_id": b.ID,
	}, bson.M{
		"$set": bson.M{
			"status": status,
			"error":  errorMessage,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Log build details
func (b *Build) Log(log string) error {

	type l struct {
		stream string `json:"stream"`
	}

	_, err := common.DB.Collection(common.BuildsCollection).UpdateOne(context.Background(), bson.M{
		"_id": b.ID,
	}, bson.M{
		"$push": bson.M{
			"logs": log,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// EnvVar struct
type EnvVar struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// BuildOptions struct
type BuildOptions struct {
	Branch string   `json:"branch" bson:"branch"`
	Env    []EnvVar `json:"env" bson:"env"`
}
