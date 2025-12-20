package database

import (
	"context"
	"distributed-configuration-management/shared/infrastructure/config"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDefault(cfg *config.Config) *mongo.Client {

	URL := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", cfg.Database.MongoDB.Username, cfg.Database.MongoDB.Password, cfg.Database.MongoDB.Host)

	if cfg.Database.MongoDB.Username == "" && cfg.Database.MongoDB.Password == "" {
		URL = fmt.Sprintf("mongodb://%s/?authSource=admin", cfg.Database.MongoDB.Host)
	}

	if cfg.Database.MongoDB.Password != "" {
		password := "xxxxxxxxxxxxxxx"
		maskedPasswordURL := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", cfg.Database.MongoDB.Username, password, cfg.Database.MongoDB.Host)
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", maskedPasswordURL)
	} else {
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", URL)
	}

	setup := options.Client()
	setup.ApplyURI(URL)
	setup.SetDirect(true)

	client, err := mongo.NewClient(setup)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("connect error >>> ", err)
		log.Fatal(err.Error())
	}

	//testing the connections
	existingCollectionNames, err := client.Database(cfg.Database.MongoDB.DbName).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println("existingCollectionNames >>> ", existingCollectionNames)

	return client
}

func NewMasterConnect(cfg *config.Config) *mongo.Client {

	URL := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", cfg.Database.MasterDB.Username, cfg.Database.MasterDB.Password, cfg.Database.MasterDB.Host)

	if cfg.Database.MasterDB.Username == "" && cfg.Database.MasterDB.Password == "" {
		URL = fmt.Sprintf("mongodb://%s/?authSource=admin", cfg.Database.MasterDB.Host)
	}

	if cfg.Database.MasterDB.Password != "" {
		password := "xxxxxxxxxxxxxxx"
		maskedPasswordURL := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=admin", cfg.Database.MasterDB.Username, password, cfg.Database.MasterDB.Host)
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", maskedPasswordURL)
	} else {
		fmt.Printf("\n>>>>>>> MongoDB URI : %s\n\n", URL)
	}

	setup := options.Client()
	setup.ApplyURI(URL)
	setup.SetDirect(true)

	client, err := mongo.NewClient(setup)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Connect(context.Background())
	if err != nil {
		fmt.Println("connect error >>> ", err)
		log.Fatal(err.Error())
	}

	//testing the connections
	existingCollectionNames, err := client.Database(cfg.Database.MasterDB.DbName).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println("existingCollectionNames >>> ", existingCollectionNames)

	return client
}

type MongoWithoutTransaction struct {
	MongoClient *mongo.Client
}

func NewMongoWithoutTransaction(c *mongo.Client) *MongoWithoutTransaction {
	return &MongoWithoutTransaction{MongoClient: c}
}

func (r *MongoWithoutTransaction) GetDatabase(ctx context.Context) (context.Context, error) {
	session, err := r.MongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)

	return sessionCtx, nil
}

func (r *MongoWithoutTransaction) Close(ctx context.Context) error {
	mongo.SessionFromContext(ctx).EndSession(ctx)
	return nil
}

//----------------------------------------------------------------------------------------

type MongoWithTransaction struct {
	MongoClient *mongo.Client
}

func NewMongoWithTransaction(c *mongo.Client) *MongoWithTransaction {
	return &MongoWithTransaction{MongoClient: c}
}

func (r *MongoWithTransaction) BeginTransaction(ctx context.Context) (context.Context, error) {

	session, err := r.MongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, session)

	err = session.StartTransaction()
	if err != nil {
		panic(err)
	}

	return sessionCtx, nil
}

func (r *MongoWithTransaction) CommitTransaction(ctx context.Context) error {

	err := mongo.SessionFromContext(ctx).CommitTransaction(ctx)
	if err != nil {
		return err
	}

	mongo.SessionFromContext(ctx).EndSession(ctx)

	return nil
}

func (r *MongoWithTransaction) RollbackTransaction(ctx context.Context) error {

	err := mongo.SessionFromContext(ctx).AbortTransaction(ctx)
	if err != nil {
		return err
	}

	mongo.SessionFromContext(ctx).EndSession(ctx)

	return nil
}

func (r *MongoWithTransaction) PrepareCollection(databaseName string, collectionNames []string) {
	db := r.MongoClient.Database(databaseName)

	existingCollectionNames, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}

	mapCollName := map[string]int{}
	for _, name := range existingCollectionNames {
		mapCollName[name] = 1
	}

	for _, name := range collectionNames {
		if _, exist := mapCollName[name]; !exist {
			// r.createCollection(db.Collection(name), db)
		}
	}
}

func (r *MongoWithTransaction) CreateIndex(dbName, collectionName, dateField string) error {

	db := r.MongoClient.Database(dbName)

	index := mongo.IndexModel{
		Keys: bson.M{dateField: 1},
	}

	_, err := db.Collection(collectionName).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoWithTransaction) CreateIndexedTTL(dbName, collectionName, dateField string, expireInSecond int32) error {

	db := r.MongoClient.Database(dbName)

	index := mongo.IndexModel{
		Keys:    bson.M{dateField: 1},
		Options: options.Index().SetExpireAfterSeconds(expireInSecond),
	}

	_, err := db.Collection(collectionName).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoWithTransaction) createCollection(coll *mongo.Collection, db *mongo.Database) {
	createCmd := bson.D{{"create", coll.Name()}}
	res := db.RunCommand(context.Background(), createCmd)
	err := res.Err()
	if err != nil {
		panic(err)
	}
}

func (r *MongoWithTransaction) SaveOrUpdate(ctx context.Context, databaseName, collectionName string, id string, data any) (any, error) {

	coll := r.MongoClient.Database(databaseName).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", data}}
	opts := options.Update().SetUpsert(true)

	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%v %v %v", result.UpsertedCount, result.ModifiedCount, result.UpsertedID), nil
}
