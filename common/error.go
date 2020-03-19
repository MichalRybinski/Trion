package common


import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
)
//changes gojsonschema result.Errors() into JSON formatted string
//Useful for APIs (as a part of error response body)
func JSONSchemaValidationErrorsToString (result *gojsonschema.Result) string {
	message := make([]string, 0)
		for _, desc := range result.Errors() {
			message = append(message, fmt.Sprintf("%v", desc))
		}
	return Lines2JSONString(&message)
}

// basic code from https://github.com/kataras/iris/blob/a1e9813c610a61aba3e803a655cefc36c6efc3b2/_examples/tutorial/mongodb/httputil/error.go

var validStackFuncs = []func(string) bool{
	func(file string) bool {
		return strings.Contains(file, "/mongodb/api/")
	},
}

// RuntimeCallerStack returns the app's `file:line` stacktrace
// to give more information about an error cause.
func RuntimeCallerStack() (s string) {
	var pcs [10]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		for _, fn := range validStackFuncs {
			if fn(frame.File) {
				s += fmt.Sprintf("\n\t\t\t%s:%d", frame.File, frame.Line)
			}
		}

		if !more {
			break
		}
	}

	return s
}

// HTTPError describes an HTTP error.
type HTTPError struct {
	error `json:"-"`
	Stack       string    `json:"-"` // the whole stacktrace.
	//TODO: wrap HTTPError in general error struct, as it's not safe to release call stack as API response e.g. via failJSON
	CallerStack string    `json:"-"` // the caller, file:lineNumber
	When        time.Time `json:"-"` // the time that the error occurred.
	// ErrorCode int: maybe a collection of known error codes.
	StatusCode int `json:"statusCode"`
	// could be named as "reason" as well
	//  it's the message of the error.
	Details *json.RawMessage `json:"details,ommitempty"`
}

func (err HTTPError) writeHeaders(ctx iris.Context) {
	ctx.StatusCode(err.StatusCode)
	ctx.Header("X-Content-Type-Options", "nosniff")
}

func newError(statusCode int, err error, format string, args ...interface{}) HTTPError {
	if format == "" {
		format = http.StatusText(statusCode)
	}
	fmt.Println(args...)
	desc := fmt.Sprintf(format, args...)
	if err == nil {
		err = errors.New(desc)
	}
	details := PrepJSONRawMsg(desc)
	
	return HTTPError{
		err,
		string(debug.Stack()),
		RuntimeCallerStack(),
		time.Now(),
		statusCode,
		details,
	}
}

// LogFailure will print out the failure to the "logger".
func LogFailure(logger io.Writer, ctx iris.Context, err HTTPError) {
	timeFmt := err.When.Format("2006/01/02 15:04:05")
	firstLine := fmt.Sprintf("%s %s: %s", timeFmt, http.StatusText(err.StatusCode), err.Error())
	whitespace := strings.Repeat(" ", len(timeFmt)+1)
	fmt.Fprintf(logger, "%s\n%sIP: %s\n%sURL: %s\n%sSource: %s\n",
		firstLine, whitespace, ctx.RemoteAddr(), whitespace, ctx.FullRequestURI(), whitespace, err.CallerStack)
}

// Fail will send the status code, write the error's reason
// and return the HTTPError for further use, i.e logging, see `InternalServerError`.
func Fail(ctx iris.Context, statusCode int, err error, format string, args ...interface{}) HTTPError {
	httpErr := newError(statusCode, err, format, args...)
	httpErr.writeHeaders(ctx)
	errD, _ := httpErr.Details.MarshalJSON()
	ctx.WriteString(string(errD))
	return httpErr
}

// FailJSON will send to the client the error data as JSON.
// Meant for API error responses.
func FailJSON(ctx iris.Context, statusCode int, err error, format string, args ...interface{}) HTTPError {
	httpErr := newError(statusCode, err, format, args...)
	errD, _ := httpErr.Details.MarshalJSON()
	fmt.Println("FailJSON, errD: ",string(errD))
	httpErr.writeHeaders(ctx)
	fmt.Println("FailJSON: before ctx.JSON")
	ctx.JSON(httpErr)
	return httpErr
}

// InternalServerError logs to the server's terminal
// and dispatches to the client the 500 Internal Server Error.
// Internal Server errors are critical, so we log them to the `os.Stderr`.
func InternalServerError(ctx iris.Context, err error, format string, args ...interface{}) {
	LogFailure(os.Stderr, ctx, Fail(ctx, iris.StatusInternalServerError, err, format, args...))
}

// InternalServerErrorJSON acts exactly like `InternalServerError` but instead it sends the data as JSON.
// Useful for APIs.
func InternalServerErrorJSON(ctx iris.Context, err error, format string, args ...interface{}) {
	LogFailure(os.Stderr, ctx, FailJSON(ctx, iris.StatusInternalServerError, err, format, args...))
}


// bad request after validating payload response helper
func BadRequestAfterJSchemaValidationResponse(ctx iris.Context, result *gojsonschema.Result) {
	var msg=""
	if result != nil { msg = JSONSchemaValidationErrorsToString(result)}
	httpErr := FailJSON(ctx,iris.StatusBadRequest,fmt.Errorf("Invalid request structure"),"%s",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}

//general bad request after some error
func BadRequestAfterErrorResponse(ctx iris.Context, err error) {
	var msg=""
	if err != nil { msg = err.Error()}
	httpErr := FailJSON(ctx,iris.StatusBadRequest,err,"%v",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}

//general not found after some error
func NotFoundAfterErrorResponse(ctx iris.Context, err error) {
	var msg=""
	if err != nil { msg = err.Error()}
	httpErr := FailJSON(ctx,iris.StatusNotFound,err,"%v",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}

//conflict 409 if another instance of resource exists
func ConflictAfterErrorResponse(ctx iris.Context, err error) {
	var msg=""
	if err != nil { msg = err.Error()}
	httpErr := FailJSON(ctx,iris.StatusConflict,err,"%v",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}
// classic 401
func UnauthorizedResponse(ctx iris.Context, err error) {
	var msg=""
	if err != nil { msg = err.Error()}
	httpErr := FailJSON(ctx,iris.StatusUnauthorized,err,"%v",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}

// classic 403
func ForbiddenResponse(ctx iris.Context, err error) {
	var msg=""
	if err != nil { msg = err.Error()}
	httpErr := FailJSON(ctx,iris.StatusForbidden,err,"%v",msg)
	LogFailure(os.Stderr, ctx, httpErr)
}


func APIErrorSwitch(ctx iris.Context, err error, additionalMsg string) {
	switch err.(type) {
		case InvalidIdError: BadRequestAfterErrorResponse(ctx,err)
		case NotFoundError: NotFoundAfterErrorResponse(ctx,err)
		case ItemAlreadyExistsError: ConflictAfterErrorResponse(ctx,err)
		case UnauthorizedError: UnauthorizedResponse(ctx,err)
		default: InternalServerErrorJSON(ctx, err, "%v", err.Error()) // general 500 error			
	}
}

type ItemAlreadyExistsError struct {
	Item string
}
func (e ItemAlreadyExistsError) Error() string { return e.Item }

type InvalidIdError struct {
	Id string
}
func (e InvalidIdError) Error() string { return e.Id +  " : invalid id"}

type NotFoundError struct {
	Item string
}
func (e NotFoundError) Error() string { return e.Item }

type NotSupportedDBError struct {
	DBType string
}
func (e NotSupportedDBError) Error() string { return e.DBType + " : not supported" }

type UnauthorizedError struct {
	Info string
}
func (e UnauthorizedError) Error() string { return e.Info }

type InvalidParametersError struct {
	Parameters map[string]interface{}
}
func (e InvalidParametersError) Error() string { 
	return MapStringInterface2String(e.Parameters)
}