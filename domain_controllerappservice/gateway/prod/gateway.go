package prod

import (
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/gateway/prod"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/database"
	"distributed-configuration-management/shared/infrastructure/logger"
)

type gateway struct {
	*database.MongoWithTransaction
	*prod.AgentImpl
	*prod.GlobalConfigImpl
	*prod.ConfigImpl
}

func NewGateway(log logger.Logger, appData driver.ApplicationData, cfg *config.Config) *gateway {

	cl := database.NewMongoDefault(cfg)
	mwt := database.NewMongoWithTransaction(cl)

	dbName := cfg.Database.MongoDB.DbName

	// prod.PrepareCollection(dbName, mwt)

	return &gateway{
		MongoWithTransaction: mwt,
		AgentImpl:            &prod.AgentImpl{MongoWithTransaction: mwt, Log: log, DbName: dbName},
		GlobalConfigImpl:     &prod.GlobalConfigImpl{MongoWithTransaction: mwt, Log: log, DbName: dbName},
		ConfigImpl:           &prod.ConfigImpl{Log: log, Cfg: cfg},
	}
}
