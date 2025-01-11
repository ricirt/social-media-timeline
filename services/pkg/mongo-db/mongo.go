package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	once   sync.Once
)

// InitMongoClient initializes MongoDB connection
func InitMongoClient(uri string) error {
	var initErr error
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientOptions := options.Client().
			ApplyURI(uri).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(10 * time.Second)

		client, initErr = mongo.Connect(ctx, clientOptions)
		if initErr != nil {
			initErr = fmt.Errorf("MongoDB connection failed: %v", initErr)
			return
		}

		if err := client.Ping(ctx, nil); err != nil {
			initErr = fmt.Errorf("MongoDB ping failed: %v", err)
			return
		}

		log.Println("Successfully connected to MongoDB")
	})

	return initErr
}

// GetMongoClient returns MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is not initialized")
	}
	return client, nil
}
