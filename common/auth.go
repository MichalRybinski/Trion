package common

import (
	//"os"
	"github.com/kataras/iris/v12"
	"github.com/iris-contrib/middleware/jwt"
	"time"
	"strconv"
)

// Prepares jwt token given specific claims
func GetJWTToken(ctx iris.Context, claims map[string]interface{}) string {
	claims["iat"]=strconv.FormatInt(time.Now().Unix(),10)
	claims["iss"]=ctx.Host()
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(TrionConfig.SecretKey))
	return tokenString
}

// reads claims from jwt token
func GetClaimsFromJWTToken(ctx iris.Context) map[string]interface{} {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	var ret = map[string]interface{}{}
	for key, value := range user.Claims.(jwt.MapClaims) {
			ret[key]=value
	}
	return ret
}

// Error handler for jwt.Config.ErrorHandler
func OnJWTError(ctx iris.Context, err error) {
	if err == nil {
		return
	}

	ctx.StopExecution()
	UnauthorizedResponse(ctx,err)
}