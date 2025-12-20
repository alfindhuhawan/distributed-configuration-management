package config

type AgentInit struct {
	Name          string `json:"name"`
	IP            string `json:"ip"`
	WorkerURL     string `json:"worker_url"`
	AgentCacheLoc string `json:"agent_cache_loc"`
}
