package getconfig

import (
	"context"
	"distributed-configuration-management/shared/pkg/workersdk"
)

//go:generate mockery --name Outport -output mocks/

type getConfigInteractor struct {
	worker workersdk.Worker
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(worker workersdk.Worker) Inport {
	return &getConfigInteractor{
		worker: worker,
	}
}

// Execute the usecase
func (r *getConfigInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	hitResp, err := r.worker.Hit(ctx)
	if err != nil {
		return nil, err
	}

	res.HitResponse = hitResp

	return res, nil
}
