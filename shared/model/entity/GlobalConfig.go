package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalConfig struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	URL                 string             `bson:"url" json:"url"`
	PollIntervalSeconds int                `bson:"poll_interval_seconds" json:"poll_interval_seconds"`
	Version             float64            `bson:"version" json:"version"`
	Type                string             `bson:"type" json:"type"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}
