package shared

import (
	"context"
	"time"

	"github.com/aureleoules/heapstate/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

// App struct
type App struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"-" bson:"user_id"`

	Token            string           `json:"token" validate:"required" bson:"token"`
	Provider         Provider         `json:"provider" bson:"provider"`
	BuildOptions     BuildOptions     `json:"build_options" validate:"required" bson:"build_options"`
	ContainerOptions ContainerOptions `json:"container_options" bson:"container_options"`
	Owner            string           `json:"owner" validate:"required" bson:"owner"`
	Name             string           `json:"name" validate:"required" bson:"name"`

	ContainerID string `json:"-" bson:"container_id"`

	CompleteURL string `json:"complete_url" bson:"complete_url"`
	URL         string `json:"url" bson:"url"`

	LastBuild Build     `json:"last_build" bson:"-"` // Dynamic
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type DeployStatus int

const (
	Building DeployStatus = iota
	Deploying
	Deployed
	Idle
	BuildError
	DeployError
)

type ContainerOptions struct {
	MaxRAM float64 `json:"max_ram" bson:"max_ram"`
}

// SetContainerID util function
func (app *App) SetContainerID(id string) error {
	_, err := common.DB.Collection(common.AppsCollection).UpdateOne(context.Background(), bson.M{
		"_id": app.ID,
	}, bson.M{
		"$set": bson.M{
			"container_id": id,
		},
	})
	return err
}

// Validate deployment
func (app *App) Validate() error {
	validate := validator.New()
	return validate.Struct(app)
}

// Save deployment
func (app *App) Save() error {
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	app.ID = primitive.NewObjectID()
	_, err := common.DB.Collection(common.AppsCollection).InsertOne(context.Background(), app)
	return err
}
