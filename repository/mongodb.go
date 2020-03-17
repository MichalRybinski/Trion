package repository

import (
	//"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"

	"github.com/MichalRybinski/Trion/common"
	"github.com/MichalRybinski/Trion/common/models"
	"fmt"
	//"encoding/json"

	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
	"time"
	"regexp"
	//"github.com/fatih/structs"
)

type mongoDBHandler struct {
	MongoClientOptions *options.ClientOptions
	MongoClient *mongo.Client
	MongoSystemDB *mongo.Database
	MongoProjectsDB *mongo.Collection
}

var MongoDBHandler mongoDBHandler

func (mh *mongoDBHandler) MongoDBInit(appConfig *common.AppConfig) {
	
	mh.MongoSystemDB = mh.MongoClient.Database(common.SysDBName)
	//for mongodb DB & collection will be created with first insert to collection
	mh.MongoProjectsDB = mh.MongoSystemDB.Collection(appConfig.DBConfig.MongoConfig.ProjectsColl)
	users, err := mh.initSystemUsers()
	fmt.Println("== init, users: ",users)
	//check if system project "trion" already exists before inserting anything
	itemsMap, err := mh.GetDocs(common.SysDBName, 
		appConfig.DBConfig.MongoConfig.ProjectsColl, 
		bson.M{"name":"trion"})
	if err != nil {
		if _, ok := err.(common.NotFoundError); !ok { log.Fatal(err) }
		fmt.Println("To był not found, len(itemsMap)",len(itemsMap))
	}
	now := time.Now()
	if len(itemsMap) == 0 {
		fmt.Println("Inserting Trion...")
		mh.InsertOne(common.SysDBName,
			appConfig.DBConfig.MongoConfig.ProjectsColl,
			bson.M{"name" : "trion", 
			"type" : "system", 
			"schema_rev": "1",
			"owner":users[0]["_id"].(string),
			"createdAt" : now,
			"updatedAt" : now,	
		})
	}
	initSysIndexes(mh.MongoProjectsDB,sysProjIndexModels)
	initSysIndexes(mh.MongoProjectsDB,userIndexModels)
}

func (mh *mongoDBHandler) initSystemUsers() ([]map[string]interface{}, error) {
	var err error
	var sysAdmin models.UserDBModel
	var now time.Time
	var hashedPassword []byte
	// check if default system admin exists
	itemsMap, err := mh.GetDocs(common.SysDBName, common.DBUsersCollectionName, bson.M{"login":"sysadmin"})
	if err !=nil { 
		if _, ok := err.(common.NotFoundError); !ok { goto Done }
	}
	if len(itemsMap) == 0 {
		//insert default admin user
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte("sysadmin"), 12)
		if err != nil {goto Done}
		now = time.Now()
		sysAdmin = models.UserDBModel{
			Login: "sysadmin",
			Hash: string(hashedPassword),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if itemsMap, err = mh.InsertOne(common.SysDBName, common.DBUsersCollectionName, sysAdmin); err != nil { goto Done }
	}
	Done:
	return itemsMap, err
}


func ConvertStringIDToObjID(stringID string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(stringID)
	if err != nil { err = common.InvalidIdError{stringID} } //Maybe wrap original error?
	return oid, err
}
// returns slice of acquired docs or error
func (mh *mongoDBHandler) GetDocs(dbname string, 
				collectionname string, 
				filter interface{}) ([]map[string]interface{}, error) {
	
	fmt.Println("=> GetDocs, filter: %v", filter)
	var err error
	var itemsMap []map[string]interface{}
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	parsedFilter = common.ConvertInterfaceToMapStringInterface(filter)
	// check if "_id" is part of filter, convert accordingly
	if _, ok := parsedFilter["_id"]; ok {
		parsedFilter["_id"], err = ConvertStringIDToObjID(parsedFilter["_id"].(string))
		if err != nil { /* log */ goto Done }
	}

	itemsMap, err = mh.getDocs(dbname,collectionname,parsedFilter)
	if err != nil {
		//log
	}
	Done:
	return itemsMap, err
}

func (mh *mongoDBHandler) InsertOne(dbname string, 
				collectionname string, 
				doc interface{}) ([]map[string]interface{}, error) {

	var itemsMap []map[string]interface{}
	var err error
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	
	var insDoc = map[string]interface{}{}
	insDoc = common.ConvertInterfaceToMapStringInterface(doc)
	// check if "_id" is part of request, convert accordingly
	if _, ok := insDoc["_id"]; ok {
		insDoc["_id"], err = ConvertStringIDToObjID(insDoc["_id"].(string))
		if err != nil { /* log */ return itemsMap, err }
	}

	var res *mongo.InsertOneResult
	if err == nil {
		now := time.Now()
		
		insDoc["createdAt"] = now
		insDoc["updatedAt"] = now
		// insDoc["owner"] = now
		res, err = collection.InsertOne(context.TODO(), insDoc)
		if err == nil {
			fmt.Printf("inserted document with ID %v\n", res.InsertedID.(primitive.ObjectID).Hex())
			itemsMap, err = mh.getDocs(dbname,collectionname,bson.M{"_id":res.InsertedID})
			fmt.Println("Inserted doc: ",itemsMap)
		} else { 
			//v, _ := err.(type)
			hasDupEntry, msgToPass := containsWriteErrDupEntry(err)
			if hasDupEntry { err = common.ItemAlreadyExistsError{msgToPass} }
		}
	}
	fmt.Printf("InsertOne: ItemsMap: %s\n InsertOne: err: %s\n",itemsMap,err) 
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
	fmt.Println("=== filter: ",filter)
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
	parsedFilter=common.ConvertInterfaceToMapStringInterface(filter)
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
		err = common.NotFoundError{ fmt.Sprintf( "not found _id : %s", parsedFilter["_id"].(primitive.ObjectID).Hex() ) }
	}
	Done:
	return itemsMap, err
}

func (mh *mongoDBHandler) UpdateDoc(dbname string, 
				collectionname string, 
				filter interface{}, 
				doc interface{}) ([]map[string]interface{}, error) {

	var itemsMap []map[string]interface{}
	var err error
	db := mh.MongoClient.Database(dbname)
	collection := db.Collection(collectionname)
	var res *mongo.UpdateResult
	//unify filter before passing to actual query
	var parsedFilter = map[string]interface{}{}
	parsedFilter=common.ConvertInterfaceToMapStringInterface(filter)
	// check if "_id" is part of filter, convert accordingly
	if _, ok := parsedFilter["_id"]; ok {
		parsedFilter["_id"], err = ConvertStringIDToObjID(parsedFilter["_id"].(string))
		if err != nil { return nil, err }
	}
	
	// parse doc into proper MongoDB update specification
	// basically: { "$set" : {doc}}
	var updateDoc = map[string]interface{}{}
	updateDoc["$set"]=doc
	updateDoc["$set"].(map[string]interface{})["updatedAt"]=time.Now()

	res, err = collection.UpdateOne(context.TODO(), parsedFilter, updateDoc)
	if err != nil {
		goto Done
	}
	if res.MatchedCount != 0 {
		itemsMap, err = mh.getDocs(dbname,collectionname,parsedFilter)
		fmt.Printf("Updated existing document %v\n for filter %v\n", itemsMap, parsedFilter)
	} else {
		err = common.NotFoundError{ fmt.Sprintf( "not found _id : %s", parsedFilter["_id"].(primitive.ObjectID).Hex() ) }
		fmt.Printf("No document updated for filter %v\n", parsedFilter)
	}
	Done:
	return itemsMap, err
}

func listExistingIndexes(coll *mongo.Collection){
	//var indexView *mongo.IndexView
	indexView := coll.Indexes()
	// Specify the MaxTime option to limit the amount of time the operation can run on the server
	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := indexView.List(context.TODO(), opts)
	if err != nil {
			log.Fatal(err)
	}

	// Get a slice of all indexes returned and print them out.
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
	}
	fmt.Println(results)
}

func initSysIndexes(coll *mongo.Collection, iModel []mongo.IndexModel) {
	indexName, err := coll.Indexes().CreateMany(
    context.Background(),
    iModel,
	)
	fmt.Println("indexName: ",indexName, " err: ",err)
}

var sysProjIndexModels = []mongo.IndexModel{
		{
				Keys: bson.D{{"name", 1},{"owner", 1}},
		},
		{
				Keys:    bson.D{{"name", 1}},
				Options: options.Index().SetUnique(true),
		},
	}

var userIndexModels = []mongo.IndexModel{
	{
			Keys:    bson.D{{"name", 1}},
			Options: options.Index().SetUnique(true),
	},
}

// private
// if err returned from mongo write operation contains duplicate entry
// e.g. breaking unique index
// returns true/false and the "{...}" part of original error message if true
func containsWriteErrDupEntry(err error) (bool, string) {
	containsDup := false
	var errMsg string
	if v, ok := err.(mongo.WriteException); ok {
		for idx, werr:=range v.WriteErrors {
			//log stuff before anything gets altered
			fmt.Println("err.WriteErrors[",idx,"].Index=",werr.Index)
			fmt.Println("err.WriteErrors[",idx,"].Code=",werr.Code)
			fmt.Println("err.WriteErrors[",idx,"].Message=",werr.Message)
			// err code 11000 or 11001 in MongoDB indicates duplicate key
			if werr.Code == 11000 || werr.Code == 11001 {
				containsDup = true
				// get the dup key msg
				pat := regexp.MustCompile(`({)(.*?)(})`)
				errMsg = pat.FindString(werr.Message)
				break;
			}
		}
	}
  return containsDup,errMsg
}