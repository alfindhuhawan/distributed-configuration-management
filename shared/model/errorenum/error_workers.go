package errorenum

const (
	ErrorInternalOnWorker  ErrorType = "ERWORKER00 %s"
	ErrorWorker            ErrorType = "ERWORKER01 worker configuration not set"
	ErrorURLWorkerRequired ErrorType = "ERWORKER02 url is required"
	ErrorConfigWorker      ErrorType = "ERWORKER03 save config for worker error : %s"
)
