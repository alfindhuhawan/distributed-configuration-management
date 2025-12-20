package registeragent

import (
	"context"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/util"
	"distributed-configuration-management/shared/model/entity"
	"distributed-configuration-management/shared/model/enum"
	"distributed-configuration-management/shared/model/errorenum"
	"os"

	"github.com/google/uuid"
)

//go:generate mockery --name Outport -output mocks/

type registerAgentInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &registerAgentInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *registerAgentInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {
	res := &InportResponse{}

	var uuidAgent string

	// validation name cant be empty
	if req.Name == "" {
		return nil, errorenum.ErrorAgentNameCantBeEmpty
	}

	// validation config available or not
	var pollURL string
	var pollIntervalSeconds int
	var lastVersion float64
	globalConfig, _ := r.outport.FindOneGlobalConfig(ctx, string(enum.DefaultGlobalConfigEnum))
	if globalConfig != nil {
		pollURL = globalConfig.URL
		pollIntervalSeconds = globalConfig.PollIntervalSeconds
		lastVersion = globalConfig.Version
	}

	uuidAgent = uuid.NewString()

	err := r.outport.SaveAgent(ctx, &entity.Agent{
		AgentID:   uuidAgent,
		Name:      req.Name,
		IP:        req.IP,
		CreatedAt: req.TimeNow,
	})
	if err != nil {
		return nil, err
	}

	// check data agent for cache start
	urlCache := r.outport.GetConfigAgentCacheLoc(ctx)
	isExist := util.FileExists(urlCache)
	if !isExist {
		// create new file if not exist
		file, err := os.Create(urlCache)
		if err != nil {
			return nil, errorenum.ErrorCreateAgentCache
		}
		defer file.Close()
	}

	// save file
	err = util.SaveCache(urlCache, config.AgentCache{
		AgentID:             uuidAgent,
		AgentName:           req.Name,
		LastVersion:         lastVersion,
		PollIntervalSeconds: pollIntervalSeconds,
		URL:                 pollURL,
		LastPollAt:          "",
	})
	if err != nil {
		return nil, errorenum.ErrorSaveAgentCache.Var(err.Error())
	}

	res.AgentID = uuidAgent
	res.PollIntervalSeconds = pollIntervalSeconds
	res.PollURL = pollURL
	res.AgentName = req.Name

	return res, nil
}
