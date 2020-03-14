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
	fmt.Println("=> prepareParsedFilter, parsedFilter: %v", parsedFilter)
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
				jsonDoc map[string]interface{}) ([]map[string]interface{}) {

	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	res, err := collection.InsertOne(context.TODO(), jsonDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID.(primitive.ObjectID).Hex())
	itemsMap, _ := mh.getDocs(dbname,collectionname,bson.M{"_id":res.InsertedID})
	
	return itemsMap
}

func (mh *mongoDBHandler) getDocs(dbname string, 
																		collectionname string, 
																		filter interface{}) ([]map[string]interface{}, error) {
	fmt.Println("== getDocs")
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	fmt.Println("=== filter: %v",filter)
	cursor, err := collection.Find(context.TODO(),filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var results []bson.M
	var itemsMap []map[string]interface{}
	
	if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
			return nil, err
	}
	if len(results) <= 0 {
		fmt.Println("No doc found")
	} else {
		fmt.Println("Doc(s) found:")
		for _, result := range results {
			var itemMap map[string]interface{}
			b, _ := bson.Marshal(result)
			bson.Unmarshal(b, &itemMap)
			itemMap["_id"] = itemMap["_id"].(primitive.ObjectID).Hex()
			fmt.Printf("itemMap after id: %v\n",itemMap)
			itemsMap = append(itemsMap, itemMap)
		}
		fmt.Printf("itemsMap: %s\n",itemsMap)
		for k, v := range itemsMap {
			fmt.Println("itemsMap[%v] v: %v", k, v)
		}
	}
	fmt.Println("== /getDocs")
	return itemsMap, nil
}

func (mh *mongoDBHandler) DeleteDoc(dbname string, 
				collectionname string, 
					filter interface{}) ([]map[string]interface{}, error) {
	
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	prepareParsedFilter(parsedFilter, filter)
	oid, err := primitive.ObjectIDFromHex(parsedFilter["_id"].(string))
	if err != nil { err = common.InvalidIdError{parsedFilter["_id"].(string)}; return nil, err }
	parsedFilter["_id"]=oid
	// grab deleted doc, so it can be provided in response
	itemsMap, err := mh.getDocs(dbname,collectionname,parsedFilter)
	if err != nil {return nil, err }
  fmt.Println("== DeleteDoc, doc to be deleted: ", itemsMap)
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	res, err := collection.DeleteOne(context.TODO(), parsedFilter)
	if err != nil {
    return nil, err
	}
	fmt.Printf("== DeleteDoc, deleted %v documents\n", res.DeletedCount)
	return itemsMap, nil
	
}

