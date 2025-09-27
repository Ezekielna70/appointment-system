package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Appointment struct {
	AppointmentID primitive.ObjectID   `bson:"_id,omitempty" json:"appointment_id,omitempty"`
	Title         string               `bson:"title" json:"title"`
	CreatorID     primitive.ObjectID   `bson:"creator_id" json:"creator_id"`
	Participants  []primitive.ObjectID `bson:"participants" json:"participants"`
	StartTime     time.Time            `bson:"start_time" json:"start_time"`
	EndTime       time.Time            `bson:"end_time" json:"end_time"`
}