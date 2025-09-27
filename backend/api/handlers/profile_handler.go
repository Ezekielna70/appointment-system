package handlers

import (
	"net/http"
	"time"

	"backend/api/models"
	"backend/config"
	"backend/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetProfile GET the currently logged-in user's profile
func GetProfile(c *gin.Context) {
	userCollection := database.GetCollection("users")
	userIDStr, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	var user models.User
	err := userCollection.FindOne(c.Request.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the currently logged-in user's profile, including Fullname and Timezone Preference
func UpdateProfile(c *gin.Context) {
	userCollection := database.GetCollection("users")
	userIDStr, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))

	var input struct {
		Name              string `json:"name"`
		PreferredTimezone string `json:"preferred_timezone"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}


	update := bson.M{
		"$set": bson.M{
			"name":               input.Name,
			"preferred_timezone": input.PreferredTimezone,
		},
	}

	_, err := userCollection.UpdateOne(c.Request.Context(), bson.M{"_id": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}


	var updatedUser models.User
	userCollection.FindOne(c.Request.Context(), bson.M{"_id": userID}).Decode(&updatedUser)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": updatedUser.ID,
		"usr": updatedUser.Username,
		"tz":  updatedUser.PreferredTimezone,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	cfg := config.LoadConfig()
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"token":   tokenString, 
	})
}