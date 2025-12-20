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

type AgentImpl struct {
	*database.MongoWithTransaction
	Log    logger.Logger
	DbName string
}

func (r *AgentImpl) SaveAgent(ctx context.Context, obj *entity.Agent) error {
	coll := r.MongoClient.Database(r.DbName).Collection(CollectionAgent)

	_, err := coll.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (r *AgentImpl) UpdateAgent(ctx context.Context, name string, obj *repository.AgentUpdateRequest) error {
	coll := r.MongoClient.Database(r.DbName).Collection(CollectionAgent)

	criteria := bson.M{}
	criteria["name"] = name

	updated := bson.M{}
	updated["$set"] = obj

	_, err := coll.UpdateOne(ctx, criteria, updated)
	if err != nil {
		return err
	}

	return nil
}

func (r *AgentImpl) FindOneAgent(ctx context.Context, name string) (*entity.Agent, error) {
	var obj entity.Agent

	coll := r.MongoClient.Database(r.DbName).Collection(CollectionAgent)

	criteria := bson.M{}
	criteria["name"] = name

	err := coll.FindOne(ctx, criteria).Decode(&obj)
	if err != nil {
		r.Log.Error(ctx, err.Error())

		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("agent not found")
		}

		return nil, err
	}

	return &obj, nil
}
