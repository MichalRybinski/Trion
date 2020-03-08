package API

import (
	"os"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/context"
	
	"github.com/xeipuuv/gojsonschema"
	"fmt"
	//"encoding/json"
	"github.com/MichalRybinski/Trion/schemas"
	//"Trion/repository"
	"github.com/MichalRybinski/Trion/common"	
)

type ProjectService struct {
	ProjectSchema *gojsonschema.Schema
}

/*type projRequest struct{
	
}
*/
func NewProjectService() *ProjectService {
	p := new(ProjectService)
	p.ProjectSchema, _ = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemas.ProjectJSchema))
	return p
}

//func prepErrRespJSON(ctx iris.Context, errdetails []string)

func (p *ProjectService) Save(ctx iris.Context) {
	var projRequest map[string]interface{}
	//var httpErr common.HTTPError
	err := ctx.ReadJSON(&projRequest)
	if err != nil {
		httpErr := common.FailJSON(ctx,iris.StatusBadRequest,err,"%v",err.Error())
		common.LogFailure(os.Stderr, ctx, httpErr)
		return 
	}
	docLoader := gojsonschema.NewGoLoader(projRequest)
	result, err := p.ProjectSchema.Validate(docLoader)
	if err != nil {
		httpErr := common.FailJSON(ctx,iris.StatusBadRequest,err,"%v",err.Error())
		common.LogFailure(os.Stderr, ctx, httpErr)
		return
	}
	if result.Valid() {
		ctx.StatusCode(iris.StatusOK)
		//TODO: make it a response with content created after save to data storage
	} else {
		  httpErr := common.FailJSON(ctx,iris.StatusBadRequest,fmt.Errorf("Invalid request payload"),"%s",common.JSONSchemaValidationErrorsToString(result))
		  common.LogFailure(os.Stderr, ctx, httpErr)
	}
}
