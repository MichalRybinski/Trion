package API
import "github.com/kataras/iris/v12"
// generic service interface
type Service interface {
	Create(ctx iris.Context)
	GetAll(ctx iris.Context)
	Delete(ctx iris.Context)
}