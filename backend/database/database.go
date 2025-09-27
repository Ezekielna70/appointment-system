package database

import (
	"context"
	"log"
	"time"
	"backend/config" 

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var cfg config.Config 


// To connect to database
func ConnectDB() {
	cfg = config.LoadConfig() 
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client
	log.Println("Connected to MongoDB!")
}

// GetCollection now uses the stored config
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Database(cfg.DBName).Collection(collectionName)
}