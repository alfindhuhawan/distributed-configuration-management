package config

type Redis struct {
	Address      		string `json:"address"`
	Port         		int    `json:"port"`
	Db           		int    `json:"db"`
	Password     		string `json:"password"`
	TimeDuration 		int    `json:"time_duration"` // IN MINUTE
	TimeDurationSetting int    `json:"time_duration_setting"` // IN MINUTE
}

type Cache struct {
	Redis Redis `json:"redis"`
}
