package getconfig

import (
	"context"
	"distributed-configuration-management/shared/model/enum"
)

//go:generate mockery --name Outport -output mocks/

type getConfigInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &getConfigInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *getConfigInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// code your usecase definition here ...
	//!
	typeConfig := req.Type
	if typeConfig == "" {
		typeConfig = string(enum.DefaultGlobalConfigEnum)
	}

	configObj, _ := r.outport.FindOneGlobalConfig(ctx, typeConfig)

	isDataChanged := false
	var respURL string
	var pollIntervalSeconds int
	var latestVersion float64
	if configObj != nil {
		respURL = configObj.URL
		pollIntervalSeconds = configObj.PollIntervalSeconds
		latestVersion = configObj.Version

		// check version first to now data changed
		if configObj.Version > req.Version {
			isDataChanged = true
		}
	}

	res.URL = respURL
	res.PollIntervalSeconds = pollIntervalSeconds
	res.Version = latestVersion
	res.IsDataChanged = isDataChanged

	return res, nil
}
