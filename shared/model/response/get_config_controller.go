package response

import "distributed-configuration-management/shared/model/repository"

type GetConfigControllerResponse struct {
	Code        int                             `json:"code"`
	CodeMessage string                          `json:"code_message"`
	Data        repository.GlobalConfigResponse `json:"data"`
}
