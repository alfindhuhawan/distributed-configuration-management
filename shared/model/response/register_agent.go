package response

import "distributed-configuration-management/shared/model/repository"

type RegisterResponse struct {
	Code        int                              `json:"code"`
	CodeMessage string                           `json:"code_message"`
	Data        repository.RegisterAgentResponse `json:"data"`
}
