package getconfig

import (
	"context"
	"distributed-configuration-management/shared/model/response"
)

// mirza here

// Inport of Usecase
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase
type InportRequest struct {
}

// InportResponse is response payload after running the usecase
type InportResponse struct {
	*response.HitResponse
}

func (r InportRequest) Validate() error {
	return nil
}
