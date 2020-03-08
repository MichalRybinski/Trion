package API
import "github.com/kataras/iris/v12"
// generic service interface
type Service interface {
	Save(ctx iris.Context)
	Get(ctx iris.Context)
	Delete(ctx iris.Context)
}