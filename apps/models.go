package apps

import (
	"context"
	"time"

	"github.com/aureleoules/heapstack/common"
	"github.com/aureleoules/heapstack/deploys"
	"github.com/aureleoules/heapstack/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
)

// App struct
type App struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"-" bson:"user_id"`

	Token        string               `json:"token" validate:"required"`
	Provider     shared.Provider      `json:"provider"`
	BuildOptions deploys.BuildOptions `json:"build_options" validate:"required"`
	Owner        string               `json:"owner" validate:"required"`
	Name         string               `json:"name" validate:"required"`

	URL string `json:"url" bson:"url"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Validate deployment
func (d *App) Validate() error {
	validate := validator.New()
	return validate.Struct(d)
}

// Save deployment
func (d *App) Save() error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	d.ID = primitive.NewObjectID()
	_, err := common.DB.Collection(common.AppsCollection).InsertOne(context.Background(), d)
	return err
}
