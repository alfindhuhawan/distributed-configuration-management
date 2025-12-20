package response

import "net/http"

type HitResponse struct {
	URL        string      `json:"url"`
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
}

type ConfigWorkerResponse struct {
	Code        int    `json:"code"`
	CodeMessage string `json:"code_message"`
}
