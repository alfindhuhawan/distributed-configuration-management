package errorenum

const (
	ErrorInternalOnAgent             ErrorType = "ERAGENT00 %s"
	ErrorEmptyAgent                  ErrorType = "ERAGENT01 please set agents first on config.json"
	ErrorRegisterAgent               ErrorType = "ERAGENT02 register agent error : %s"
	ErrorCreateAgentCache            ErrorType = "ERAGENT03 error create agent cache file"
	ErrorGetConfigController         ErrorType = "ERAGENT04 get config from controller error : %s"
	ErrorGetConfigControllerNoChange ErrorType = "ERAGENT05 %s"
	ErrorSaveAgentCache              ErrorType = "ERAGENT06 error save agent cache file : %s"
	ErrorAgentNameCantBeEmpty        ErrorType = "ERAGENT07 agent name cant be empty"
)
