package repository

import (
	"context"
	//"errors"
	//"log"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	// up to you:
	//"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	//"encoding/json"

	c "github.com/MichalRybinski/Trion/common"
)

type DataStoreService interface {
	Create(ctx context.Context, createReq map[string]interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error)
	Update(ctx context.Context, filter interface{}, updateReq map[string]interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error)
	Read(ctx context.Context, filter interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error)
	Delete(ctx context.Context, filter interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error)
	IsTokenNotRecalled(ctx context.Context, uid string, uuid string) (bool, error)
}

type dataStoreService struct {
	dbtype string
}

var _ DataStoreService = (*dataStoreService)(nil)

func NewDataStoreService(dbType string) (DataStoreService, error) {
	var dss dataStoreService
	var err error
	switch dbType {
		// add psql in future here - with fallthrough
		case "mongodb": dss = dataStoreService{dbtype: dbType}
		case "default" : err = c.NotSupportedDBError{dbType}
	}
	return &dss, err
}


func (dss *dataStoreService)Create(ctx context.Context, createReq map[string]interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error
	fmt.Println("DSS Create...")
	if dbData["dbName"].(string) != "" && dbData["collectionName"].(string) != "" {
		switch dss.dbtype {
			case "mongodb": {
					itemsMap, err = MongoDBHandler.InsertOne(dbData["dbName"].(string), 
						dbData["collectionName"].(string),
						createReq)
						fmt.Println("DSS: InsertOne err %T %v",err,err)
					if err !=nil { /*log?*/}
			}
			default:
		}
	} else {
		err = c.InvalidParametersError{dbData}
	}
	fmt.Printf("DSS Create, itemsMap: %v err: %v \n",itemsMap, err)
	return itemsMap, err
}


func (dss *dataStoreService)Read(ctx context.Context, filter interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error
	fmt.Println("DSS Read...")
	if dbData["dbName"].(string) != "" && dbData["collectionName"].(string) != "" {
		switch dss.dbtype {
			case "mongodb": {
					itemsMap, err = MongoDBHandler.GetDocs(dbData["dbName"].(string), 
					dbData["collectionName"].(string),
					filter)
					if err != nil {/* log? */}
			}
			default:
		}
	}	else {
		err = c.InvalidParametersError{dbData}
	}

	return itemsMap, err
}

func (dss *dataStoreService)Update(ctx context.Context, filter interface{}, updateReq map[string]interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error
	fmt.Println("DSS Update...")
	if dbData["dbName"].(string) != "" && dbData["collectionName"].(string) != "" {
		switch dss.dbtype {
			case "mongodb": {
					itemsMap, err = MongoDBHandler.UpdateDoc(dbData["dbName"].(string), 
					dbData["collectionName"].(string),
					filter,
					updateReq)
					if err != nil {/* log? */}
			}
			default:
		}
	}	else {
		err = c.InvalidParametersError{dbData}
	}

	return itemsMap, err
}

func (dss *dataStoreService)Delete(ctx context.Context, filter interface{}, dbData map[string]interface{}) ([]map[string]interface{}, error) {
	var itemsMap []map[string]interface{}
	var err error
	fmt.Println("DSS Delete...")
	if dbData["dbName"].(string) != "" && dbData["collectionName"].(string) != "" {
		switch dss.dbtype {
			case "mongodb": {
					itemsMap, err = MongoDBHandler.DeleteDoc(dbData["dbName"].(string), 
					dbData["collectionName"].(string),
					filter)
					if err != nil {/* log? */}
			}
			default:
		}
	}	else {
		err = c.InvalidParametersError{dbData}
	}

	return itemsMap, err
}

func (dss *dataStoreService)IsTokenNotRecalled(ctx context.Context, uid string, uuid string) (bool, error) {
	var itemsMap []map[string]interface{}
	var err error
	var res = false
	switch dss.dbtype {
		case "mongodb": {
			filter := map[string]interface{}{
				"userId" : uid,
				"uuid" : uuid,
			}
			itemsMap, err = MongoDBHandler.GetDocs(c.UsersDBName, 
				c.UsersDBAuthCollection,
				filter)
			if err != nil {/* log? */}
		}
		default:
	}
	if err == nil && len(itemsMap) > 0 { 
		/*log*/ res = true
	} else if err == nil && len(itemsMap) < 1 {
		err = c.UnauthorizedError{"Token is recalled"}
	}
	return res, err
}