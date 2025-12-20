package updateconfig

import "distributed-configuration-management/shared/model/repository"

// Outport of usecase
type Outport interface {
	repository.FindOneGlobalConfigRepo
	repository.SaveConfigRepo
	repository.UpdateConfigRepo
}
