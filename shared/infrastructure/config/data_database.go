package config

type MongoDB struct {
	DbName   string `json:"database"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MasterDB struct {
	DbName   string `json:"database"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Database struct {
	MongoDB  MongoDB  `json:"mongodb"`
	MasterDB MasterDB `json:"master_db"`
}
