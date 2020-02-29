package common

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holder
var DB *mongo.Database

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
	DB = db
	return db
}

// GetDatabase returns db
func GetDatabase() *mongo.Database {
	return DB
}
