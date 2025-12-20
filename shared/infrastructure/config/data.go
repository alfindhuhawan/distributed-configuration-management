package config

type Config struct {
	ApplicationServer   ApplicationServer   `json:"application_server"`
	ApplicationInternal ApplicationInternal `json:"application_internal"`
	Credentials         Credentials         `json:"credentials"`
	PollIntervalSeconds int                 `json:"poll_interval_seconds"`
	BaseDelaySeconds    int                 `json:"based_delay_seconds"`
	MaxDelaySeconds     int                 `json:"max_delay_seconds"`
	JWTSecretKey        string              `json:"jwt_secret_key"`
	Database            Database            `json:"database"`
	AgentInit           AgentInit           `json:"agent_init"`
	Worker              Worker              `json:"worker"`
}

type AgentCache struct {
	AgentID             string  `json:"agent_id"`
	AgentName           string  `json:"agent_name"`
	LastVersion         float64 `json:"last_version"`
	PollIntervalSeconds int     `json:"poll_interval_seconds"`
	URL                 string  `json:"url"`
	LastPollAt          string  `json:"last_poll_at"`
}
