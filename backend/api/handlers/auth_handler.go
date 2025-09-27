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
	"go.mongodb.org/mongo-driver/mongo"
)

// Login function for logging in User
func Login(c *gin.Context) {
	userCollection := database.GetCollection("users")

	var body struct {
		Username string `json:"username"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User
	err := userCollection.FindOne(c.Request.Context(), bson.M{"username": body.Username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User does not exist, tell the frontend to go to signup option
			c.JSON(http.StatusNotFound, gin.H{
                "error": "User not found", 
                "action": "redirect_to_signup",
            })
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// If login is possible, generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"usr": user.Username,
		"tz":  user.PreferredTimezone,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	cfg := config.LoadConfig()
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Signup creates a new user and logs them in
func Signup(c *gin.Context) {
    userCollection := database.GetCollection("users")

    var input struct {
        Name              string `json:"name"`
        Username          string `json:"username"`
        PreferredTimezone string `json:"preferred_timezone"`
    }

    if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

    // Check if username inserted already exists
    count, _ := userCollection.CountDocuments(c.Request.Context(), bson.M{"username": input.Username})
    if count > 0 {
        c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
        return
    }

    newUser := models.User{
        ID: primitive.NewObjectID(),
        Name: input.Name,
        Username: input.Username,
        PreferredTimezone: input.PreferredTimezone,
    }

    _, err := userCollection.InsertOne(c.Request.Context(), newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": newUser.ID,
		"usr": newUser.Username,
		"tz":  newUser.PreferredTimezone,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	cfg := config.LoadConfig()
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
}