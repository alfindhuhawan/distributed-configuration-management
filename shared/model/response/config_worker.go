package response

type HitResponse struct {
	URL              string `json:"url"`
	StatusCode       int    `json:"status_code"`
	ExternalResponse any    `json:"response"`
}

type ConfigWorkerResponse struct {
	Code        int    `json:"code"`
	CodeMessage string `json:"code_message"`
}
