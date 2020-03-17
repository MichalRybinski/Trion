package API

import (
	//"os"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/v12/hero"
	//"github.com/kataras/iris/context"
	
	//"github.com/xeipuuv/gojsonschema"
	"fmt"
	//"strings"

	//"github.com/MichalRybinski/Trion/schemas"
	"github.com/MichalRybinski/Trion/repository"
	"github.com/MichalRybinski/Trion/common"
	"golang.org/x/crypto/bcrypt"

)

type UserService struct {
	DSS repository.DataStoreService
}

func NewUserService(dss repository.DataStoreService) *UserService {
	us := new(UserService)
	us.DSS = dss
	return us
}

func (us *UserService) SignIn(ctx iris.Context, projectName string) {
	var dbData = map[string]interface{}{ "collectionName" : common.DBUsersCollectionName }
  if projectName == "" { 
		dbData["dbName"] = common.SysDBName
	} else {
		dbData["dbName"] = projectName
	}

	fmt.Println("== UserService, Signin, dbData: ",dbData)
	var request map[string]interface{}
	err := ctx.ReadJSON(&request)
	if err != nil {
		common.BadRequestAfterErrorResponse(ctx,err)
		return
	}
	// pass on login/pwd and grab result, response based on it
	// get matching user from DB
	var filter = map[string]interface{}{} //init empty
	if _, ok := request["login"].(string); ok {
		filter["login"] = request["login"].(string)
	}

	itemsMap, err := us.DSS.Read(nil, filter, dbData); 
	if err == nil {
		if len(itemsMap) !=1 { // if not just one matching entry...
			goto Unauthorized
		} else {
			if err = bcrypt.CompareHashAndPassword([]byte(itemsMap[0]["hash"].(string)),
				[]byte(request["password"].(string))); err == nil {
					common.StatusJSON(ctx,iris.StatusOK,"%s","Authorized")
					return
			}	else { // If the two passwords don't match, return a 401 status
				goto Unauthorized 
			}
		}
	} else if  _, ok := err.(common.NotFoundError); !ok { // if no entry found - unauthorized, else regular
		common.APIErrorSwitch(ctx,err,"Error during sign-in")
	}
	Unauthorized:
	common.UnauthorizedResponse(ctx,common.UnauthorizedError{"Credentials"})
	return
}