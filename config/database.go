package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Client *mongo.Client

// ConnectDB establishes a connection to MongoDB
func ConnectDB() error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://pratama:cjzMmK1k7BGQ8zhq@yoginara.vvsjt.mongodb.net/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Test connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	Client = client
	log.Println("Connected to MongoDB!")
	return nil
}

// GetCollection returns a reference to a MongoDB collection
func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("stepupDB").Collection(collectionName)
}

// DisconnectDB closes the MongoDB connection
func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Disconnect(ctx); err != nil {
		log.Println("Error disconnecting from MongoDB:", err)
		return
	}
	log.Println("Disconnected from MongoDB")
}
