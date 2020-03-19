package API

import (
	"github.com/MichalRybinski/Trion/repository"
	c "github.com/MichalRybinski/Trion/common"
	"github.com/kataras/iris/v12"
)

// A middleware for processing JWT Tokens and verifying authorization

type AuthHandler struct {
	DSS repository.DataStoreService
	AuthDB map[string]interface{}
}

func NewAuthHandler(dss repository.DataStoreService) *AuthHandler {
	ah := new(AuthHandler)	
	ah.DSS = dss
	ah.AuthDB = map[string]interface{}{
		"dbName" : c.UsersDBName,
		"collectionName" : c.UsersDBAuthCollection,
	}
	return ah
}

func (ah *AuthHandler) AllowAccess(ctx iris.Context) {
	var ok bool
	var err error

	if ok, err = ah.Authenticate(ctx); ok {
		if ok, err = ah.Authorize(); ok { 
			ctx.Next() //exit point
		} else { 
			c.ForbiddenResponse(ctx,err)
		}
	} else {
		c.UnauthorizedResponse(ctx,err)
	}
	ctx.StopExecution() //sorry, no access
}

func (ah *AuthHandler) Authenticate(ctx iris.Context) (bool, error) {
	var authenticated = false
	claims:= c.GetClaimsFromJWTToken(ctx)
	ok, err := ah.isTokenNotRecalled(claims["sub"].(string), claims["uuid"].(string))
	if err == nil && ok { authenticated = true }
	return authenticated, err
}

func (ah *AuthHandler) Authorize() (bool, error) {
	var err error
	// consider authorizing via usage of policies, see either casbin for very simple parsing or
	// https://github.com/ory/ladon for AWS-alike policies and more flexibility
	// maybe deploy Ory Keto as a separate policy handing service - still
	// Ory Keto requires it's own DB, CockroachDB?
	return true, err
}

// Verifying whether the valid token is still present in Auths or not
// TODO put Auths stuff into Redis than in DB in future - will be faster
func (ah *AuthHandler) isTokenNotRecalled(uid string, uuid string) (bool, error) {
	var notRecalled = false
	filter := map[string]interface{}{
		"uid" : uid,
		"uuid" : uuid,
	}
	
	itemsMap, err := ah.DSS.Read(nil, filter, ah.AuthDB)
	if err == nil && len(itemsMap) == 1 { //>1 multiple entries? Something went really south, leave result false
		/*log*/ notRecalled = true
	} else if err == nil && len(itemsMap) < 1 {
		err = c.UnauthorizedError{"Token is recalled"}
	}

	return notRecalled, err
}