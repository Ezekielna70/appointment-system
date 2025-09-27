package main

import (
	"context"
	"log"
	"strings"

	"backend/api/handlers"
	"backend/api/middleware"
	"backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	database.ConnectDB()
	seedUsers() 

	router := gin.Default()

	// Public routes
	router.POST("/login", handlers.Login)
	router.POST("/signup", handlers.Signup) 

	// Serve frontend files
	router.Static("/public", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// Protected routes(access only with token)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", handlers.GetProfile)
		api.PUT("/profile", handlers.UpdateProfile)
		api.GET("/users", handlers.GetUsers)
		api.POST("/appointments", handlers.CreateAppointment)
		api.GET("/appointments", handlers.GetAppointments)
	}

	log.Println("Server starting on port 8080...")
	router.Run(":8080")
}

// Add user dummy
func seedUsers() {
	userCollection := database.GetCollection("users")

	// Create unique index
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		if !strings.Contains(err.Error(), "Index already exists") {
			log.Printf("Could not create index: %v", err)
		}
	}

	users := []interface{}{
		bson.M{"name": "Alice", "username": "alice", "preferred_timezone": "Asia/Jakarta"},
		bson.M{"name": "Bob", "username": "bob", "preferred_timezone": "Europe/London"},
		bson.M{"name": "Charlie", "username": "charlie", "preferred_timezone": "America/New_York"},
	}

	// Loop to add users
	for _, user := range users {
		username := user.(bson.M)["username"]
		count, err := userCollection.CountDocuments(context.TODO(), bson.M{"username": username})
		if err == nil && count == 0 {
			_, err := userCollection.InsertOne(context.TODO(), user)
			if err != nil {
				log.Printf("Could not insert user %s: %v", username, err)
			}
		}
	}
	log.Println("User seeding complete.")
}