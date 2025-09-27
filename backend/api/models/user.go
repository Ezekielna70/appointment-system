package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name"`
	Username          string             `bson:"username" json:"username"`
	PreferredTimezone string             `bson:"preferred_timezone" json:"preferred_timezone"` // e.g., "Asia/Jakarta", "Pacific/Auckland"
}