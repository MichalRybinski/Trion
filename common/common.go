package common

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/kataras/iris/v12"
)

func ReadJsonToString(filepath string) (jsonString string) {
	jsonFile, err := os.Open(filepath)
	// if we os.Open returns an error then handle it
	if err != nil {
    	fmt.Println(err)
	}
	fmt.Println("Successfully opened %s+\n", filepath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	jsonString = string(byteValue)
	fmt.Println(jsonString)
	return jsonString
}

// make sure 'desc' won't cause rawMessage marshalling errors
// we want always an object or table of strings in 'details'
func PrepJSONRawMsg(desc string) *json.RawMessage {
	var details json.RawMessage
	if (strings.HasPrefix(desc, "{") || strings.HasPrefix(desc, "[") ) {
		details = json.RawMessage(desc)
	} else {
		msg, _ := json.Marshal([]string {desc})
		details = json.RawMessage(string(msg))
	}
	return &details
}

// common action among Service implementations
// if request JSON can not be unmarshalled, returns 400 via iris.Context
func ParseRequestToJSON(ctx iris.Context, projRequest *map[string]interface{}) error {
	err := ctx.ReadJSON(projRequest)
	if err != nil {
		BadRequestAfterErrorResponse(ctx,err)
		return err
	}
	return nil
}