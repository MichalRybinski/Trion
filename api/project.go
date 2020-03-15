package API

import (
	//"os"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/v12/hero"
	//"github.com/kataras/iris/context"
	
	"github.com/xeipuuv/gojsonschema"
	//"fmt"
	"strings"

	"github.com/MichalRybinski/Trion/schemas"
	"github.com/MichalRybinski/Trion/repository"
	"github.com/MichalRybinski/Trion/common"

)


type ProjectService struct {
	ProjectSchema *gojsonschema.Schema
	ProjectFilterSchema *gojsonschema.Schema
	PSS repository.ProjectStoreService
}

func NewProjectService(pss repository.ProjectStoreService) *ProjectService {
	p := new(ProjectService)
	p.ProjectSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemas.ProjectJSchema))
	p.ProjectFilterSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemas.ProjectFilterJSchema))
	p.PSS = pss
	return p
}

// Saves a new project to the database after POST call to interface
func (p *ProjectService) Create(ctx iris.Context) {

	var projRequest map[string]interface{}
	err := common.ParseRequestToJSON(ctx, &projRequest)
	if err != nil { return }
	
	docLoader := gojsonschema.NewGoLoader(projRequest)
	result, err := p.ProjectSchema.Validate(docLoader)
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return
	}

	if result.Valid() {
		projRequest["schema_rev"] = schemas.ProjectJSchemaVersion
		itemsMap, err:= p.PSS.Create(nil, projRequest) //save to DB
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
		if err != nil {
			common.BadRequestAfterErrorResponse(ctx,err)
			return
		}

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
	itemsMap, err:= p.PSS.Read(nil, filter); 
	if err !=nil {
		common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}

func (p *ProjectService) GetById(ctx iris.Context, id string) {
	var filter = map[string]interface{}{ "_id" : id }
	// get projects from DB
	itemsMap, err:= p.PSS.Read(nil, filter); 
	if err !=nil {
		common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}


func (p *ProjectService) DeleteById(ctx iris.Context, id string) {
	
	var filter = map[string]interface{}{ "_id" : id }
	itemsMap, err := p.PSS.Delete(nil, filter)

	if err != nil {
		common.APIErrorSwitch(ctx,err,"")
	} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	}
	return
}

func (p* ProjectService) UpdateById(ctx iris.Context, id string) {
	var projRequest map[string]interface{}
	err := common.ParseRequestToJSON(ctx, &projRequest)
	if err != nil { common.BadRequestAfterErrorResponse(ctx,err); return }

	var filter = map[string]interface{}{ "_id" : id }

	docLoader := gojsonschema.NewGoLoader(projRequest)
	result, err := p.ProjectSchema.Validate(docLoader)
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return
	}

	if result.Valid() {
		projRequest["schema_rev"] = schemas.ProjectJSchemaVersion
		itemsMap, err:= p.PSS.Update(nil, filter, projRequest) //save to DB
		if err != nil { 
			common.APIErrorSwitch(ctx,err,common.SliceMapToJSONString(itemsMap))
		} else {
		common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
		}
	} else { common.BadRequestAfterJSchemaValidationResponse(ctx,result) }
	return

}

