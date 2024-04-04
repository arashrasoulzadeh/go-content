package providers

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoconnection() *mongo.Database {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("database uri not provided")
	}
	db := os.Getenv("MONGODB_NAME")
	if db == "" {
		panic("database name not provided")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client.Database(db)
}
