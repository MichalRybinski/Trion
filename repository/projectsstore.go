package repository

import (
	"context"
	//"errors"
	//"log"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// up to you:
	//"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	//"encoding/json"

	"github.com/MichalRybinski/Trion/common"
)

type ProjectStoreService interface {
	Create(ctx context.Context, createReq map[string]interface{}) ([]map[string]interface{}, error)
	Update(ctx context.Context, filter interface{}, updateReq map[string]interface{}) ([]map[string]interface{}, error)
	Read(ctx context.Context, filter interface{}) ([]map[string]interface{}, error)
	Delete(ctx context.Context, filter interface{}) ([]map[string]interface{}, error)
}

type projectStoreService struct {
	dbtype string
	mongoC *mongo.Collection
	//add psql in future here
}

var _ ProjectStoreService = (*projectStoreService)(nil)

func NewProjectStoreService(collection interface{}) ProjectStoreService {
	var pss projectStoreService
	t, ok := collection.(*mongo.Collection)
	if ok ==true {
		pss = projectStoreService{dbtype: "mongodb", mongoC: t}
		fmt.Println("\npss: %v",pss)
  } else {
	// add psql in future here
		fmt.Println("Unhandled collection type")
	}
	return &pss
}

func (pss *projectStoreService) Create(ctx context.Context, createReq map[string]interface{}) ([]map[string]interface{}, error) {
	
	var itemsMap []map[string]interface{}
	var err error
	
	switch pss.dbtype {
		case "mongodb" : {
			projName := createReq["name"].(string)
			fmt.Println("PROJSTORESERVICE.Create %s",projName)
			//TODO: insert doc into system-Projects if doc doesn't exist or return error with existing doc
			itemsMap, err = MongoDBHandler.GetDocs("TrionSystem", 
				common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
				map[string]interface{}{"name": projName})
			if err != nil {
				return nil, err
			} else if len(itemsMap) > 0 {
        err = common.ProjectAlreadyExistsError{projName}
			} else { //len(itemsMap) < 1
				itemsMap, err = MongoDBHandler.InsertDoc("TrionSystem",
					common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
					createReq)
				fmt.Printf("PSS, itemsMap: %v err: %v \n",itemsMap, err)
				if err !=nil { return itemsMap, err }
			}
		}
	//TODO: in case of psql create separate DB with name equal to "name" in project request - use repository.proper handler
		default:
			fmt.Println("Error, unhandled db type")
	}
	return itemsMap, err
}

func (pss *projectStoreService) Read(ctx context.Context, filter interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error
	
	switch pss.dbtype {
		case "mongodb": {
			itemsMap, err = MongoDBHandler.GetDocs("TrionSystem", 
				common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
				filter)
			if err != nil {
				// log
			}
		}
		default:
	}

	return itemsMap, err
}

func (pss *projectStoreService) Delete(ctx context.Context, filter interface{}) ([]map[string]interface{}, error) {
	
	var itemsMap []map[string]interface{}
	var err error
	
	switch pss.dbtype {
		case "mongodb": {
			itemsMap, err = MongoDBHandler.DeleteDoc("TrionSystem", 
				common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
				filter)
			if err != nil {
				//log
			}
		}
		default:
	}

	return itemsMap, err
}

func (pss *projectStoreService) Update(ctx context.Context, filter interface{}, updateReq map[string]interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error

	switch pss.dbtype {
		case "mongodb": {
			itemsMap, err = MongoDBHandler.UpdateDoc("TrionSystem", 
				common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
				filter,
				updateReq)
			if err != nil {
				//log
			}
		}
		default:
	}

	return itemsMap, err
}