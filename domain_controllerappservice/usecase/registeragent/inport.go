package registeragent

import (
	"context"
	"distributed-configuration-management/shared/model/repository"
)

// mirza here

// Inport of Usecase
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase
type InportRequest struct {
	repository.RegisterAgentRequest
}

// InportResponse is response payload after running the usecase
type InportResponse struct {
	repository.RegisterAgentResponse
}

func (r InportRequest) Validate() error {
	return nil
}
