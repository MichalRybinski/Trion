package API

import (
	"os"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/context"
	
	"github.com/xeipuuv/gojsonschema"
	"fmt"
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
	err := ctx.ReadJSON(&projRequest)
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return 
	}
	
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
			var httpErr common.HTTPError
			// if project with 'name' already exists, return 409 along with existing resource data
			if projExists, ok := err.(repository.ProjectAlreadyExistsError); ok {
				httpErr = common.FailJSON(ctx,iris.StatusConflict,projExists,"%v",common.SliceMapToJSONString(itemsMap))
			} else { //500 for anything else - db comms failed in general
				httpErr = common.FailJSON(ctx,iris.StatusInternalServerError,err,"%v",err.Error())
			}
			common.LogFailure(os.Stderr, ctx, httpErr)
			return
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
	{
		"filter" : {our filter json}
	}
	*/
	fmt.Println("=> API.GetAll content-type requested: %s", ctx.GetContentTypeRequested())
	if strings.ToLower(ctx.GetContentTypeRequested()) == "application/json" {
		var request map[string]interface{}
		err := ctx.ReadJSON(&request)
		if err != nil {
			common.BadRequestAfterErrorResponse(ctx,err)
			return
		}
		fmt.Println("=> API.GetAll, request: %v",request)
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
	fmt.Println("=> API.GetAll, filter: %v",filter)
	// get projects from DB
	itemsMap, err:= p.PSS.Get(nil, filter); 
	if err !=nil {
		common.InternalServerErrorJSON(ctx, err, "%v", err.Error())
	}
	common.StatusJSON(ctx,iris.StatusOK,"%v",common.SliceMapToJSONString(itemsMap))
	return
}
