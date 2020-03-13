package common

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	"github.com/kataras/iris/v12"
)

type HTTPResponse struct {
	When        time.Time `json:"-"` // the time that the error occurred.
	StatusCode int `json:"statusCode"`
	//  it's the message of the response.
	Details *json.RawMessage `json:"details"`
}

func (resp HTTPResponse) writeHeaders(ctx iris.Context) {
	ctx.StatusCode(resp.StatusCode)
	ctx.Header("X-Content-Type-Options", "nosniff")
}

func newResponse(statusCode int, format string, args ...interface{}) HTTPResponse {
	if format == "" {
		format = http.StatusText(statusCode)
	}

	desc := fmt.Sprintf(format, args...)
	details := PrepJSONRawMsg(desc)
	
	return HTTPResponse{
		time.Now(),
		statusCode,
		details,
	}
}

func SliceMapToJSONString(itemsMap []map[string]interface{} ) string {
	j, _:= json.Marshal(itemsMap)
	res := string(j)
	return res
}

// StatusJSON will send to the client the response data as JSON.
// Meant for API non-error responses.
func StatusJSON(ctx iris.Context, statusCode int, format string, args ...interface{}) HTTPResponse {
	httpResp := newResponse(statusCode, format, args...)
	httpResp.writeHeaders(ctx)
	ctx.JSON(httpResp)
	return httpResp
}

