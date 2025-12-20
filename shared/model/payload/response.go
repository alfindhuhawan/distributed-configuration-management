package payload

import (
	"distributed-configuration-management/shared/model/errorenum"
)

type Response struct {
	Code         int    `json:"code"`
	ErrorCode    string `json:"error_code"`
	CodeMessage  string `json:"code_message"`
	Message      string `json:"message"`
	ErrorMapping any    `json:"error_mapping"`
	Data         any    `json:"data"`
}

type ResponseSubmit struct {
	Code             int    `json:"code"`
	ErrorCode        string `json:"error_code"`
	CodeMessage      string `json:"code_message"`
	Message          string `json:"message"`
	ErrorMapping     any    `json:"error_mapping"`
	Data             any    `json:"data"`
	ShowErrorMessage bool   `json:"show_error_message"`
}

type ErrorMappingRes struct {
	ErrorCodes   string `json:"error_codes"`
	ErrorMessage string `json:"error_message"`
	ErrorFrom    string `json:"error_from"`
}

func NewSuccessResponse(data any, errorMapping any) any {
	var res Response
	res.Code = 200
	res.ErrorCode = ""
	res.CodeMessage = "Success"
	res.ErrorMapping = errorMapping
	res.Message = ""
	res.Data = data

	return res
}

func NoChangeResponse(data any, errorMapping any) any {
	var res Response
	res.Code = 304
	res.ErrorCode = ""
	res.CodeMessage = "No Changes"
	res.ErrorMapping = errorMapping
	res.Message = ""
	res.Data = data

	return res
}

func NewErrorResponse(err error, traceID string) any {
	var res Response
	res.CodeMessage = "false"

	et, _ := err.(errorenum.ErrorType)
	// if !ok {
	res.Code = 500
	res.CodeMessage = err.Error()
	// res.Message = err.Error()
	// return res
	// }

	res.ErrorCode = et.Code()
	// res.CodeMessage = et.Error()
	return res
}
