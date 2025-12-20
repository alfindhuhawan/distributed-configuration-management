package prod

// func PrepareCollection(dbName string, mwt *database.MongoWithTransaction) {
// 	collectionNames := []string{
// 		CollectionLogEvent,
// 		CollectionPermission,
// 	}

// 	mwt.PrepareCollection(dbName, collectionNames)

// 	err := mwt.CreateIndex(dbName, CollectionLogEvent, "trace_id")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	err = mwt.CreateIndexedTTL(dbName, CollectionLogEvent, "deleted_at", 10)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// }
