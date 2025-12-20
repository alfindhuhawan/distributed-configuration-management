package repository

import (
	"context"
	"distributed-configuration-management/shared/model/entity"
	"time"
)

type SaveConfigRepo interface {
	SaveConfig(ctx context.Context, obj *entity.GlobalConfig) error
}

type FindOneGlobalConfigRepo interface {
	FindOneGlobalConfig(ctx context.Context, typeConfig string) (*entity.GlobalConfig, error)
}

type UpdateConfigRepo interface {
	UpdateConfig(ctx context.Context, typeConfig string, obj *GlobalConfigUpdateRequest) error
}

type GlobalConfigRequest struct {
	URL                 string    `json:"url"`
	Type                string    `json:"type"`
	PollIntervalSeconds int       `json:"poll_interval_seconds"`
	TimeNow             time.Time `json:"updated_at"`
}

type GlobalConfigUpdateRequest struct {
	URL                 string    `json:"url" bson:"url"`
	Version             float64   `json:"version" bson:"version"`
	PollIntervalSeconds int       `json:"poll_interval_seconds" bson:"poll_interval_seconds"`
	UpdatedAt           time.Time `json:"updated_at" bson:"updated_at"`
}

type GlobalConfigResponse struct {
	URL                 string  `json:"url"`
	Type                string  `json:"-"`
	Version             float64 `json:"version"`
	PollIntervalSeconds int     `json:"poll_interval_seconds"`
	IsDataChanged       bool    `json:"-"`
}
