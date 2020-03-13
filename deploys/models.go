package deploys

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Build struct
type Build struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Options BuildOptions `json:"build_options" bson:"build_options"`
	Owner   string       `json:"owner" validate:"required"`
	Name    string       `json:"name" validate:"required"`

	URL string `json:"url" bson:"url"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type BuildOptions struct {
	Branch string `json:"branch" bson:"branch"`
}
