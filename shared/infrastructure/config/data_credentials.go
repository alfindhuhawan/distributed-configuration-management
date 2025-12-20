package config

type Credentials struct {
	AdminCredential AdminCredential `json:"admin"`
	AgentCredential AgentCredential `json:"agent"`
}

type AdminCredential struct {
	SecretKey string `json:"secret_key"`
}

type AgentCredential struct {
	SecretKey string `json:"secret_key"`
}
