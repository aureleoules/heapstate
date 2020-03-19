package users

import (
	"context"
	"time"

	"github.com/aureleoules/heapstate/common"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	Credits uint64 `json:"credits" bson:"credits"`

	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Profile struct
type Profile struct {
	ID        primitive.ObjectID `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Credits   uint64             `json:"credits"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedAt time.Time          `json:"created_at"`
}

// Save to db
func (u *User) Save() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	u.ID = primitive.NewObjectID()
	_, err := common.DB.Collection(common.UsersCollection).InsertOne(context.Background(), u)
	return err
}

// Validate user
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// HashPassword hash user's password
func (u *User) HashPassword() {
	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	u.Password = string(password)
}

// Public returns public data
func (u *User) Public() Profile {
	return Profile{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		UpdatedAt: u.UpdatedAt,
		Credits:   u.Credits,
		CreatedAt: u.CreatedAt,
	}
}
