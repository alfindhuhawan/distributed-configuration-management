package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	AgentID   string             `bson:"agent_id" json:"agent_id"`
	Name      string             `bson:"name" json:"name"`
	IP        string             `bson:"ip" json:"ip"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type AgentCache struct {
	AgentID             string  `json:"agent_id"`
	AgentName           string  `json:"agent_name"`
	LastVersion         float64 `json:"last_version"`
	PollIntervalSeconds float64 `json:"poll_interval_seconds"`
	URL                 string  `json:"url"`
	LastPollAt          string  `json:"last_poll_at"`
}
