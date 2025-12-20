package updateconfig

import (
	"context"
	"distributed-configuration-management/shared/model/entity"
	"distributed-configuration-management/shared/model/enum"
	"distributed-configuration-management/shared/model/errorenum"
	"distributed-configuration-management/shared/model/repository"
)

//go:generate mockery --name Outport -output mocks/

type updateConfigInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &updateConfigInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *updateConfigInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	res := &InportResponse{}

	// version update
	var latestVersion float64 = 1

	// get config with type global
	globalConfig, _ := r.outport.FindOneGlobalConfig(ctx, string(enum.DefaultGlobalConfigEnum))
	if globalConfig == nil {
		// save new config with type global
		err := r.outport.SaveConfig(ctx, &entity.GlobalConfig{
			URL:                 req.URL,
			Type:                string(enum.DefaultGlobalConfigEnum),
			UpdatedAt:           req.TimeNow,
			PollIntervalSeconds: req.PollIntervalSeconds,
			Version:             1,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// validation check is there any changes on url or poll_interval_seconds, if not no need to increase version
		if req.URL == globalConfig.URL && req.PollIntervalSeconds == globalConfig.PollIntervalSeconds {
			return nil, errorenum.ErrorControllerNochange
		}

		// update config with type global
		err := r.outport.UpdateConfig(ctx, string(enum.DefaultGlobalConfigEnum), &repository.GlobalConfigUpdateRequest{
			URL:                 req.URL,
			Version:             globalConfig.Version + 1,
			PollIntervalSeconds: req.PollIntervalSeconds,
			UpdatedAt:           req.TimeNow,
		})
		if err != nil {
			return nil, err
		}

		latestVersion = globalConfig.Version + 1
	}

	res.Type = string(enum.DefaultGlobalConfigEnum)
	res.URL = req.URL
	res.Version = latestVersion
	res.PollIntervalSeconds = req.PollIntervalSeconds

	return res, nil
}
