package repository

import (
	"context"
	"distributed-configuration-management/shared/model/entity"
	"time"
)

type SaveAgentRepo interface {
	SaveAgent(ctx context.Context, obj *entity.Agent) error
}

type UpdateAgentRepo interface {
	UpdateAgent(ctx context.Context, name string, obj *AgentUpdateRequest) error
}

type FindOneAgentRepo interface {
	FindOneAgent(ctx context.Context, typeAgent string) (*entity.Agent, error)
}

type RegisterAgentRequest struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Type    string `json:"type"`
	TimeNow time.Time
}

type RegisterAgentResponse struct {
	AgentID             string `json:"agent_id"`
	PollURL             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
	AgentName           string `json:"agent_name"`
}

type AgentUpdateRequest struct {
	IP        string    `json:"ip" bson:"ip"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
