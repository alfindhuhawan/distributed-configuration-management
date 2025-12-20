package prod

import (
	"context"
	"distributed-configuration-management/shared/infrastructure/database"
	"distributed-configuration-management/shared/infrastructure/logger"
	"distributed-configuration-management/shared/model/entity"
	"distributed-configuration-management/shared/model/repository"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type GlobalConfigImpl struct {
	*database.MongoWithTransaction
	Log    logger.Logger
	DbName string
}

func (r *GlobalConfigImpl) SaveConfig(ctx context.Context, obj *entity.GlobalConfig) error {
	coll := r.MongoClient.Database(r.DbName).Collection(CollectionGlobalConfig)

	_, err := coll.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (r *GlobalConfigImpl) UpdateConfig(ctx context.Context, typeConfig string, obj *repository.GlobalConfigUpdateRequest) error {
	coll := r.MongoClient.Database(r.DbName).Collection(CollectionGlobalConfig)

	criteria := bson.M{}
	criteria["type"] = typeConfig

	updated := bson.M{}
	updated["$set"] = obj

	_, err := coll.UpdateOne(ctx, criteria, updated)
	if err != nil {
		return err
	}

	return nil
}

func (r *GlobalConfigImpl) FindOneGlobalConfig(ctx context.Context, typeConfig string) (*entity.GlobalConfig, error) {
	var obj entity.GlobalConfig

	coll := r.MongoClient.Database(r.DbName).Collection(CollectionGlobalConfig)

	criteria := bson.M{}
	criteria["type"] = typeConfig

	err := coll.FindOne(ctx, criteria).Decode(&obj)
	if err != nil {
		r.Log.Error(ctx, err.Error())

		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("global config not found")
		}

		return nil, err
	}

	return &obj, nil
}
