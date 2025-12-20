package repository

import (
	"context"
)

type GetConfigPollIntervalSecondsRepo interface {
	GetConfigPollIntervalSeconds(ctx context.Context) int
}

type GetConfigAgentCacheLocRepo interface {
	GetConfigAgentCacheLoc(ctx context.Context) string
}
