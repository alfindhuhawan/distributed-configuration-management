package getconfig

import "distributed-configuration-management/shared/model/repository"

// Outport of usecase
type Outport interface {
	repository.FindOneGlobalConfigRepo
}
