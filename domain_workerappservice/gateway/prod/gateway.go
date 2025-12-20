package prod

import (
	"distributed-configuration-management/shared/driver"
	"distributed-configuration-management/shared/infrastructure/config"
	"distributed-configuration-management/shared/infrastructure/database"
	"distributed-configuration-management/shared/infrastructure/logger"
)

type gateway struct {
	*database.MongoWithTransaction
}

func NewGateway(log logger.Logger, appData driver.ApplicationData, cfg *config.Config) *gateway {

	cl := database.NewMongoDefault(cfg)
	mwt := database.NewMongoWithTransaction(cl)

	// prod.PrepareCollection(dbName, mwt)

	return &gateway{
		MongoWithTransaction: mwt,
	}
}
