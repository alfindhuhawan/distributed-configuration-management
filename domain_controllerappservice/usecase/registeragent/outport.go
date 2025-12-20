package registeragent

import "distributed-configuration-management/shared/model/repository"

// Outport of usecase
type Outport interface {
	repository.SaveAgentRepo
	repository.GetConfigPollIntervalSecondsRepo
	repository.FindOneGlobalConfigRepo
	repository.FindOneAgentRepo
	repository.UpdateAgentRepo
	repository.GetConfigAgentCacheLocRepo
}
