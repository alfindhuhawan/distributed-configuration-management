package runconfig

import (
	"context"
	"distributed-configuration-management/shared/model/request"
	"distributed-configuration-management/shared/pkg/workersdk"
)

//go:generate mockery --name Outport -output mocks/

type runConfigInteractor struct {
	worker workersdk.Worker
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(worker workersdk.Worker) Inport {
	return &runConfigInteractor{
		worker: worker,
	}
}

// Execute the usecase
func (r *runConfigInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	err := r.worker.UpdateConfig(ctx, request.ConfigWorker{
		URL: req.URL,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
