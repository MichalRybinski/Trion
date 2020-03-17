package API

import (
	//"os"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/v12/hero"
	//"github.com/kataras/iris/context"
	
	"github.com/xeipuuv/gojsonschema"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/MichalRybinski/Trion/schemas"
	"github.com/MichalRybinski/Trion/repository"
	"github.com/MichalRybinski/Trion/common"

)


type ProjectService struct {
	ProjectSchema *gojsonschema.Schema
	ProjectFilterSchema *gojsonschema.Schema
	DSS repository.DataStoreService
}

func NewProjectService(dss repository.DataStoreService) *ProjectService {
	p := new(ProjectService)
	p.ProjectSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemas.ProjectJSchema))
	p.ProjectFilterSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemas.ProjectFilterJSchema))
	p.DSS = dss
	return p
}

// Saves a new project to the database after POST call to interface
func (p *ProjectService) Create(ctx iris.Context) {

	var request map[string]interface{}
	err := ctx.ReadJSON(&request)
	if err != nil { common.BadRequestAfterErrorResponse(ctx,err); return }
	
	docLoader := gojsonschema.NewGoLoader(request)
	result, err := p.ProjectSchema.Validate(docLoader)
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return
	}

	var dbData = map[string]interface{}{ 
		"dbName" : common.SysDBName, 
		"collectionName" : common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
	}

	if result.Valid() {
		request["schema_rev"] = schemas.ProjectJSchemaVersion
		itemsMap, err:= p.DSS.Create(nil, request, dbData) //save to DB
		if err != nil { 
			common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
		} else {
			common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
		}
	} else { common.BadRequestAfterJSchemaValidationResponse(ctx,result) }
	return
}

// GET: Gets all projects there (save type:system) and allows to
// put filter as payload if header contains content-type: application/json
func (p *ProjectService) GetAll(ctx iris.Context) {
	var dbData = map[string]interface{}{ 
		"dbName" : common.SysDBName, 
		"collectionName" : common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
	}
	var filter = map[string]interface{}{} //init empty
  /*
	Filtering for content-type : application/json only.
	1. Detect if header declares content-type application/json
	2. If yes - parse json. If not - don't, empty filter.
	request payload:
	{
		"filter" : {our filter json}
	}
	*/
	if strings.ToLower(ctx.GetContentTypeRequested()) == "application/json" {
		var request map[string]interface{}
		err := ctx.ReadJSON(&request)
		if err != nil { common.BadRequestAfterErrorResponse(ctx,err); return }

		//validate schema
		docLoader := gojsonschema.NewGoLoader(request)
		result, err := p.ProjectFilterSchema.Validate(docLoader)
		if err != nil {
			common.BadRequestAfterErrorResponse(ctx,err)
			return
		}
		if result.Valid() {
			v, ok := request["filter"].(map[string]interface{})
			if ok { filter = v } 
		} else { 
			common.BadRequestAfterJSchemaValidationResponse(ctx,result) 
			return
		}
	}
	// get projects from DB
	itemsMap, err:= p.DSS.Read(nil, filter, dbData); 
	err=checkIfNotFound(err, len(itemsMap), filter) 
	if err !=nil {
		common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}

func (p *ProjectService) GetById(ctx iris.Context, id string) {
	var filter = map[string]interface{}{ "_id" : id }
	var dbData = map[string]interface{}{ 
		"dbName" : common.SysDBName, 
		"collectionName" : common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
	}
	// get projects from DB
	itemsMap, err:= p.DSS.Read(nil, filter,dbData); 
	err=checkIfNotFound(err, len(itemsMap), filter) 
	if err !=nil {
		common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}


func (p *ProjectService) DeleteById(ctx iris.Context, id string) {
	
	var filter = map[string]interface{}{ "_id" : id }
	var dbData = map[string]interface{}{ 
		"dbName" : common.SysDBName, 
		"collectionName" : common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
	}
	itemsMap, err := p.DSS.Delete(nil, filter, dbData)

	if err != nil {
		common.APIErrorSwitch(ctx,err,"")
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}

func (p* ProjectService) UpdateById(ctx iris.Context, id string) {
	var request map[string]interface{}
	err := ctx.ReadJSON(&request)
	if err != nil { common.BadRequestAfterErrorResponse(ctx,err); return }

	var filter = map[string]interface{}{ "_id" : id }

	docLoader := gojsonschema.NewGoLoader(request)
	result, err := p.ProjectSchema.Validate(docLoader) //project schema doesn't allow to update id or timestamps
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return
	}
	var dbData = map[string]interface{}{ 
		"dbName" : common.SysDBName, 
		"collectionName" : common.ThisAppConfig.DBConfig.MongoConfig.ProjectsColl,
	}

	if result.Valid() {
		request["schema_rev"] = schemas.ProjectJSchemaVersion
		itemsMap, err:= p.DSS.Update(nil, filter, request, dbData) //save to DB
		if err != nil { 
			common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
		} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
		}
	} else { common.BadRequestAfterJSchemaValidationResponse(ctx,result) }
	return

}

// if there's no error, but count <=0 - returns NotFoundError
func checkIfNotFound(err error, count int, details interface{}) error {
	if err == nil && count <= 0 { 
		jsonD, _ := json.Marshal(details)
		errorMsg := fmt.Sprintf("no match : %v", string(jsonD))
		err = common.NotFoundError{errorMsg}
	}
	return err
}
