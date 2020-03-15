package repository

import (
	//"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"

	"github.com/MichalRybinski/Trion/common"
	"fmt"
	"encoding/json"

	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBHandler struct {
	MongoClientOptions *options.ClientOptions
	MongoClient *mongo.Client
	MongoSystemDB *mongo.Database
	MongoProjectsDB *mongo.Collection
}

var MongoDBHandler mongoDBHandler

func (mh *mongoDBHandler) MongoDBInit(appConfig *common.AppConfig) {
	
	mh.MongoSystemDB = mh.MongoClient.Database("TrionSystem")
	//for mongodb DB & collection will be created with first insert to collection
	mh.MongoProjectsDB = mh.MongoSystemDB.Collection(appConfig.DBConfig.MongoConfig.ProjectsColl)
	//check if system project "trion" already exists before inserting anything
	itemsMap, err := mh.getDocs("TrionSystem", 
		appConfig.DBConfig.MongoConfig.ProjectsColl, 
		bson.M{"name":"trion"})
	if err != nil {
		log.Fatal(err)
	} else if len(itemsMap) == 0 {
		mh.InsertDoc("TrionSystem",
			appConfig.DBConfig.MongoConfig.ProjectsColl,
			bson.M{"name" : "trion", "type" : "system", "schema_rev": "1"})
	}
}

//handle filter types; parsedFilter will be modified
func prepareParsedFilter(parsedFilter map[string]interface{}, filter interface{}) {
	
	switch v:=filter.(type) {
		case []byte: {
			if err := json.Unmarshal(v,&parsedFilter); err!=nil {
			}
		}
		case map[string]interface{}: {
			for key,val := range v {
				parsedFilter[key]=val
			}
		}
		case bson.M : {
			var temporaryBytes []byte
			var err error
			temporaryBytes, err = bson.MarshalExtJSON(v, true, true)
			if err == nil {
				err = json.Unmarshal(temporaryBytes, &parsedFilter)
				if err != nil {}
			}
		}
		default: //nothing, empty filter
	}
	fmt.Println("=> prepareParsedFilter, parsedFilter: ", parsedFilter)
	return
}

// returns slice of acquired docs or error
func (mh *mongoDBHandler) GetDocs(dbname string, 
				collectionname string, 
				filter interface{}) ([]map[string]interface{}, error) {
	
	fmt.Println("=> GetDocs, filter: %v", filter)
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	prepareParsedFilter(parsedFilter, filter)
	//fmt.Println("=> GetDocs, parsedFilter: %v", parsedFilter)

	itemsMap, err :=mh.getDocs(dbname,collectionname,parsedFilter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return itemsMap, nil
}

func (mh *mongoDBHandler) InsertDoc(dbname string, 
				collectionname string, 
				jsonDoc map[string]interface{}) ([]map[string]interface{}, error) {

	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	res, err := collection.InsertOne(context.TODO(), jsonDoc)
	if err != nil {
		return nil, err
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID.(primitive.ObjectID).Hex())
	itemsMap, err := mh.getDocs(dbname,collectionname,bson.M{"_id":res.InsertedID})
	
	return itemsMap, err
}

func (mh *mongoDBHandler) getDocs(dbname string, 
																		collectionname string, 
																		filter interface{}) ([]map[string]interface{}, error) {
	fmt.Println("== getDocs")
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	var result bson.M
	var results []bson.M
	var itemsMap []map[string]interface{}
	fmt.Println("=== filter: %v",filter)
	cursor, err := collection.Find(context.TODO(),filter)
	if err != nil {
		goto Done
	}
	
	if err = cursor.All(context.TODO(), &results); err != nil {
			goto Done
	}
	if len(results) <= 0 {
		fmt.Println("No doc found")
	} else {
		fmt.Println("Doc(s) found:")
		for _, result = range results {
			var itemMap map[string]interface{}
			b, _ := bson.Marshal(result)
			bson.Unmarshal(b, &itemMap)
			itemMap["_id"] = itemMap["_id"].(primitive.ObjectID).Hex()
			fmt.Printf("itemMap after id: %v\n",itemMap)
			itemsMap = append(itemsMap, itemMap)
		}	
	}	
	Done:
	fmt.Printf("itemsMap: %v\n",itemsMap)
	for k, v := range itemsMap {
		fmt.Println("itemsMap[",k,"]=",v)
	}
	fmt.Println("== /getDocs")
	return itemsMap, err
}

func ConvertStringIDToObjID(stringID string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(stringID)
	if err != nil { err = common.InvalidIdError{stringID} } //Maybe wrap original error?
	return oid, err
}

func (mh *mongoDBHandler) DeleteDoc(dbname string, 
				collectionname string, 
					filter interface{}) ([]map[string]interface{}, error) {
	
	var itemsMap []map[string]interface{}
	var err error
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	var res *mongo.DeleteResult
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	prepareParsedFilter(parsedFilter, filter)
	// check if "_id" is part of filter, convert accordingly
	if _, ok := parsedFilter["_id"]; ok {
		parsedFilter["_id"], err = ConvertStringIDToObjID(parsedFilter["_id"].(string))
		if err != nil { goto Done }
	}
	
	// grab doc to be deleted, so it can be provided in response for reference
	itemsMap, err = mh.getDocs(dbname,collectionname,parsedFilter)
	if err != nil { goto Done }
	fmt.Println("== DeleteDoc, doc to be deleted: ", itemsMap)
	res, err = collection.DeleteOne(context.TODO(), parsedFilter)
	fmt.Printf("== DeleteDoc, deleted %v documents\n", res.DeletedCount)
	if res.DeletedCount == 0 {
		err = common.NotFoundError{ fmt.Sprintf( "not found with _id : %s", parsedFilter["_id"].(primitive.ObjectID).Hex() ) }
	}
	Done:
	return itemsMap, err
}

func (mh *mongoDBHandler) UpdateDoc(dbname string, 
				collectionname string, 
				filter interface{}, 
				jsonDoc map[string]interface{}) ([]map[string]interface{}, error) {

	var itemsMap []map[string]interface{}
	var err error
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	var res *mongo.UpdateResult
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	prepareParsedFilter(parsedFilter, filter)
	// check if "_id" is part of filter, convert accordingly
	if _, ok := parsedFilter["_id"]; ok {
		parsedFilter["_id"], err = ConvertStringIDToObjID(parsedFilter["_id"].(string))
		if err != nil { return nil, err }
	}
	
	// parse jsonDoc into proper MongoDB update specification
	// basically: { "$set" : {jsonDoc}}
	var updateDoc = map[string]interface{}{}
	updateDoc["$set"]=jsonDoc

	res, err = collection.UpdateOne(context.TODO(), parsedFilter, updateDoc)
	if err != nil {
		goto Done
	}
	if res.MatchedCount != 0 {
		itemsMap, err = mh.getDocs(dbname,collectionname,parsedFilter)
		fmt.Printf("Updated existing document %v\n for filter %v\n", itemsMap, parsedFilter)
	} else {
		err = common.NotFoundError{ fmt.Sprintf( "not found with _id : %s", parsedFilter["_id"].(primitive.ObjectID).Hex() ) }
		fmt.Printf("No document updated for filter %v\n", parsedFilter)
	}
Done:
	return itemsMap, err
}