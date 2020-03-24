package common

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holder
var DB *mongo.Database

const (
	UsersCollection  = "users"
	AppsCollection   = "apps"
	BuildsCollection = "builds"
)

// InitDB database connection
func InitDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	if db == nil {
		panic(errors.New("Couldn't connect to MongoDB"))
	}
	fmt.Println("INITIALIZED MONGODB")
	DB = db
	return db
}

// GetDatabase returns db
func GetDatabase() *mongo.Database {
	return DB
}
