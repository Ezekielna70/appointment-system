package handlers

import (
	"net/http"
	"backend/api/models"
    "backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// API endpoint to GET user's data
func GetUsers(c *gin.Context) {
    userCollection := database.GetCollection("users")
	cursor, err := userCollection.Find(c.Request.Context(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var users []models.User
	if err = cursor.All(c.Request.Context(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	c.JSON(http.StatusOK, users)
}