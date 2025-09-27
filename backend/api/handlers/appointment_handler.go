package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
	"backend/api/models"
	"backend/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// For Creating Appointment
func CreateAppointment(c *gin.Context) {
	appointmentCollection := database.GetCollection("appointments")
	userCollection := database.GetCollection("users")

	var input struct {
		Title           string   `json:"title"`
		ParticipantIDs  []string `json:"participant_ids"`
		StartTime       string   `json:"start_time"`
		DurationMinutes int      `json:"duration_minutes"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	creatorIDStr, _ := c.Get("userID")
	creatorID, _ := primitive.ObjectIDFromHex(creatorIDStr.(string))

	startTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format. Use RFC3339 (UTC)."})
		return
	}

	endTime := startTime.Add(time.Duration(input.DurationMinutes) * time.Minute)

	inviteeObjectIDs := []primitive.ObjectID{}
	for _, idStr := range input.ParticipantIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID format"})
			return
		}
		inviteeObjectIDs = append(inviteeObjectIDs, id)
	}

	// Only perform working-hours validation if there are invitees, so it ignores the current user's working hours rule
	if len(inviteeObjectIDs) > 0 {
		cursor, err := userCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": inviteeObjectIDs}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch participants"})
			return
		}
		var invitees []models.User
		if err = cursor.All(context.Background(), &invitees); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not decode participants"})
			return
		}

		for _, invitee := range invitees {
			loc, err := time.LoadLocation(invitee.PreferredTimezone)
			if err != nil {
				log.Printf("Could not load location for user %s: %s", invitee.Username, invitee.PreferredTimezone)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error with timezone location"})
				return
			}

			localStartTime := startTime.In(loc)
			localHour := localStartTime.Hour()

			if localHour < 8 || localHour >= 17 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment is outside of working hours (08:00-17:00) for user: " + invitee.Username})
				return
			}
		}
	}

	allParticipantIDs := append(inviteeObjectIDs, creatorID)

	appointment := models.Appointment{
		Title:        input.Title,
		CreatorID:    creatorID,
		Participants: allParticipantIDs,
		StartTime:    startTime,
		EndTime:      endTime,
	}

	res, err := appointmentCollection.InsertOne(c.Request.Context(), appointment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Appointment created successfully!", "id": res.InsertedID})
}

// GetAppointments function with the time 
func GetAppointments(c *gin.Context) {
	appointmentCollection := database.GetCollection("appointments")
	userIDStr, _ := c.Get("userID")
	userID, _ := primitive.ObjectIDFromHex(userIDStr.(string))
	userTimezone, _ := c.Get("userTimezone")

	loc, err := time.LoadLocation(userTimezone.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user timezone"})
		return
	}

	filter := bson.M{
		"participants": bson.M{"$in": []primitive.ObjectID{userID}},
		"start_time":   bson.M{"$gte": time.Now().Add(-1 * time.Minute)},
	}

	cursor, err := appointmentCollection.Find(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	var appointments []models.Appointment
	if err = cursor.All(c.Request.Context(), &appointments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode appointments"})
		return
	}

	type AppointmentResponse struct {
		AppointmentID primitive.ObjectID `json:"appointment_id"`
		Title         string             `json:"title"`
		StartTime     string             `json:"start_time_local"`
		EndTime       string             `json:"end_time_local"`
	}

	var response []AppointmentResponse
	for _, appt := range appointments {
		response = append(response, AppointmentResponse{
			AppointmentID: appt.AppointmentID,
			Title:         appt.Title,
			StartTime:     appt.StartTime.In(loc).Format(time.RFC1123),
			EndTime:       appt.EndTime.In(loc).Format(time.RFC1123),
		})
	}

	c.JSON(http.StatusOK, response)
}