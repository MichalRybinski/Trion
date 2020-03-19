package main

import (
	//"fmt"
	//"time"

	//mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/iris-contrib/middleware/jwt"
	//"github.com/kataras/iris/context"
	//"github.com/kataras/iris/middleware/logger"
	//"github.com/kataras/iris/middleware/recover"
	//"Trion/repository"
	"github.com/MichalRybinski/Trion/API"
	c "github.com/MichalRybinski/Trion/common"
	"github.com/MichalRybinski/Trion/repository"
	"context"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const notFoundHTML = "<h1> custom http error page </h1>"

func registerErrors(app *iris.Application) {
	// set a custom 404 handler
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.HTML(notFoundHTML)
	})
}

func registerApiRoutes(app *iris.Application) {
	dss, err := repository.NewDataStoreService(c.TrionConfig.DBConfig.DBType)
	if err != nil {
		return
	}
	j:= initJWTMiddleware()
	apiMiddleware := func(ctx iris.Context) {
		ctx.Next()
	}
	// party is just a group of routes with the same prefix
	// and middleware, i.e: "/api" and apiMiddleware.
	api := app.Party("/api", apiMiddleware)
	{ // braces are optional of course, it's just a style of code
		auth := API.NewAuthHandler(dss)
		v1 := api.Party("/v1")
		projects := v1.Party("/projects")
		{
			psHandler := API.NewProjectService(dss)
			projects.Get("/", psHandler.GetAll)
			projects.Post("/", psHandler.Create)
			// hero handlers for path parameters
			projects.Delete("/{id:string}", j.Serve, auth.AllowAccess, hero.Handler(psHandler.DeleteById))
			projects.Put("/{id:string}", hero.Handler(psHandler.UpdateById))
			projects.Get("/{id:string}", hero.Handler(psHandler.GetById))
		}
		usHandler := API.NewUserService(dss)
		v1.Post("/{projName:string}/signin", hero.Handler(usHandler.SignIn))
		v1.Post("/signin", hero.Handler(usHandler.SignIn))
		v1.Post("/{projName:string}/signout", j.Serve, auth.AllowAccess, hero.Handler(usHandler.SignOut))
		v1.Post("/signout", j.Serve, auth.AllowAccess, hero.Handler(usHandler.SignOut))
	}
}

func initJWTMiddleware() *jwt.Middleware {
		j := jwt.New(jwt.Config{
				ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
						return []byte(c.TrionConfig.SecretKey), nil
				},
				SigningMethod: jwt.SigningMethodHS256,
				ErrorHandler: c.OnJWTError,
		})
		return j
}


func registerSubdomains(app *iris.Application) {
	mysubdomain := app.Party("mysubdomain.")
	// http://mysubdomain.myhost.com
	mysubdomain.Get("/", h)

	//willdcardSubdomain := app.Party("*.")
	//willdcardSubdomain.Get("/", h)
	//willdcardSubdomain.Party("/party").Get("/", h)
}

func newApp() *iris.Application {
	app := iris.New()
	registerErrors(app)
	registerApiRoutes(app)
	registerSubdomains(app)

	//app.Handle("GET", "/healthcheck", h)

	return app
}

// generic handler from example. 
// Once all references to h removed, remove it as well
func h(ctx iris.Context) {
	method := ctx.Method()       // the http method requested a server's resource.
	subdomain := ctx.Subdomain() // the subdomain, if any.

	// the request path (without scheme and host).
	path := ctx.Path()
	// how to get all parameters, if we don't know
	// the names:
	paramsLen := ctx.Params().Len()

	ctx.Params().Visit(func(name string, value string) {
		ctx.Writef("%s = %s\n", name, value)
	})
	ctx.Writef("Info\n\n")
	ctx.Writef("Method: %s\nSubdomain: %s\nPath: %s\nParameters length: %d", method, subdomain, path, paramsLen)
}


func main() {
	switch c.TrionConfig.DBConfig.DBType {
		case "mongodb": {
			var err error
			//initiate client & connection, connection pool handled by driver
			//keep it active until program is terminated
			repository.MongoDBHandler.MongoClientOptions = options.Client().ApplyURI(c.TrionConfig.DBConfig.MongoConfig.URL)
			repository.MongoDBHandler.MongoClient, err = mongo.Connect(context.Background(), 
				repository.MongoDBHandler.MongoClientOptions)
			if err != nil {
				log.Fatal(err)
			}
			err = repository.MongoDBHandler.MongoClient.Ping(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			defer repository.MongoDBHandler.MongoClient.Disconnect(context.TODO())
			repository.MongoDBHandler.MongoDBInit(c.TrionConfig)
		}
		default:
	}
	app := newApp()
	app.Logger().SetLevel("debug")
	app.Run(iris.Addr(":" + c.TrionConfig.ServerConfig.PORT), 
		iris.WithoutServerError(iris.ErrServerClosed))
}
