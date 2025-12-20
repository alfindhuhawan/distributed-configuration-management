package prod

import (
	"context"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/logger"
)

type ConfigImpl struct {
	Log logger.Logger
	Cfg *config.Config
}

func (r *ConfigImpl) GetConfigPollIntervalSeconds(ctx context.Context) int {
	pollIntervalSeconds := r.Cfg.PollIntervalSeconds

	// set default pollintervalseconds
	if pollIntervalSeconds == 0 {
		pollIntervalSeconds = 30
	}

	return pollIntervalSeconds
}

func (r *ConfigImpl) GetConfigAgentCacheLoc(ctx context.Context) string {
	agentCacheLoc := r.Cfg.AgentInit.AgentCacheLoc

	return agentCacheLoc
}
