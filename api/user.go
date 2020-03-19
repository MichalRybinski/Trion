package API

import (
	//"os"
	"github.com/kataras/iris/v12"
	//"github.com/iris-contrib/middleware/jwt"
	//"github.com/kataras/iris/v12/hero"
	//"github.com/kataras/iris/context"
	
	//"github.com/xeipuuv/gojsonschema"
	"fmt"
	//"strings"

	//"github.com/MichalRybinski/Trion/schemas"
	"github.com/MichalRybinski/Trion/repository"
	c "github.com/MichalRybinski/Trion/common"
	m "github.com/MichalRybinski/Trion/common/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"

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
	var dbData = map[string]interface{}{ 
		"dbName" : c.UsersDBName,
		"collectionName" : c.UsersDBUsersCollection, 
	}

	fmt.Println("== UserService, Signin, dbData: ",dbData)
	var request map[string]interface{}
	err := ctx.ReadJSON(&request)
	if err != nil {
		c.BadRequestAfterErrorResponse(ctx,err)
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
			err = c.UnauthorizedError{"Could not sign in. Check credentials."}
		} else {
			if err = bcrypt.CompareHashAndPassword([]byte(itemsMap[0]["hash"].(string)),
				[]byte(request["password"].(string))); err == nil {
					//issue jwtToken for API access, long expiry etc
					//will need: user id, 
					var claims = map[string]interface{}{}
					claims["sub"] = itemsMap[0]["_id"].(string)
					claims["aud"] = projectName
					// repeat until error free or retry count elapsed
					var token string
					var retries = 3
					for {
						claims["uuid"] = uuid.NewV4().String()
						token = c.GetJWTToken(ctx, claims)
						auth, _ := c.StructToMap(m.NewAuth(claims["sub"].(string), claims["uuid"].(string), token), "json")
						_, err = us.DSS.Create(nil, auth, map[string]interface{}{
							"dbName": c.UsersDBName,
							"collectionName" : c.UsersDBAuthCollection,
						})
						retries--
						fmt.Println(" -- retries: ",retries)
						if err == nil || retries == 0 { break }
						// Consider: add some sleep between retries?
					}
					/*critsection:
					claims["uuid"] = uuid.NewV4().String()
					token:= c.GetJWTToken(ctx, claims)
					auth, _ := c.StructToMap(m.NewAuth(claims["sub"].(string), claims["uuid"].(string), token), "json")
					if _, err = us.DSS.Create(nil, auth, map[string]interface{}{
							"dbName": c.UsersDBName,
							"collectionName" : c.UsersDBAuthCollection,
						}); err != nil {goto critsection} //entry in DB must match the token!
						// TODO: break after a couple of retries*/
					if err == nil { 
						c.StatusJSON(ctx,iris.StatusOK,"%s",c.MapToJSON(map[string]interface{}{
							"token" : token,
						}))
					}
			}	else { // If the two passwords don't match, return a 401 status
				err = c.UnauthorizedError{"Could not sign in. Check credentials."}
			}
		}
	} else {
		if  _, ok := err.(c.NotFoundError); ok { 
		// if no entry found - unauthorized, else regular error handling
			err = c.UnauthorizedError{"Could not sign in. Check credentials."}
		}
	}
	if err != nil { c.APIErrorSwitch(ctx,err,"Error during sign-in") }
	return
}

func (us *UserService) SignOut(ctx iris.Context, projectName string) {
	var dbData = map[string]interface{}{ 
		"dbName" : c.UsersDBName,
		"collectionName" : c.UsersDBAuthCollection, 
	}
	fmt.Println("== UserService, SignOut, dbData: ",dbData)

	// grab uid and uuid from claims and remove appropriate Auths entry
	// thus invalidating any further usage of this particular token
	claims:= c.GetClaimsFromJWTToken(ctx)
	var filter = map[string]interface{}{
		"uid": claims["sub"].(string),
		"uuid"	:	claims["uuid"].(string),
	}

	_, err := us.DSS.Delete(nil, filter, dbData); 
	if err != nil { c.APIErrorSwitch(ctx,err,"Error during sign-out") }
	return
}