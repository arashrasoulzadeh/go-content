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
	db := os.Getenv("MONGODB_NAME")
	if db == "" {
		panic("database not provided")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	return client.Database(db)
}
