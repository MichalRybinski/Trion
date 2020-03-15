package API
import "github.com/kataras/iris/v12"
// generic APIservice interface
type APIService interface {
	Create(ctx iris.Context)
	GetAll(ctx iris.Context)
	DeleteById(ctx iris.Context, id string)
	UpdateById(ctx iris.Context, id string)
	GetById(ctx iris.Context, id string)
}
