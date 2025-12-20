package payload

import "distributed-configuration-management/shared/driver"

type Args struct {
	Type      string                 `json:"type"`
	Data      any                    `json:"data"`
	Publisher driver.ApplicationData `json:"publisher"`
	TraceID   string                 `json:"traceId"`
}

type Reply struct {
	Success      bool
	ErrorMessage string
	Data         any
}
