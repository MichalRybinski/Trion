package common

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"strings"
	//"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
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
	//returning always a table of strings containing error messages
	if (strings.HasPrefix(desc, "[") ) {
		details = json.RawMessage(desc)
	} else {
		msg, _ := json.Marshal([]string {desc})
		details = json.RawMessage(string(msg))
	}
	return &details
}

// Type Conversion helper
// handle 'input' types; 'converted' will be modified
func ConvertInterfaceToMapStringInterface(input interface{}) map[string]interface{} {
	var converted = map[string]interface{}{}
	switch v:=input.(type) {
		case []byte: {
			if err := json.Unmarshal(v,&converted); err!=nil {
			}
		}
		case map[string]interface{}: {
			for key,val := range v {
				converted[key]=val
			}
		}
		case bson.M : {
			var temporaryBytes []byte
			var err error
			temporaryBytes, err = bson.MarshalExtJSON(v, true, true)
			if err == nil {
				err = json.Unmarshal(temporaryBytes, &converted)
				if err != nil {}
			}
		}
		default: //nothing, empty filter
	}
	fmt.Println("=> ConvertInterfaceToMapStringInterface, converted: ", converted)
	return converted
}

func MapStringInterface2String(mapVar map[string]interface{}) string {
	msglist := make([]string,0)
	for k,v:= range mapVar {
		msglist=append(msglist,fmt.Sprintf("Parameter '%s' has value: '%s'",k,v))
	}
	return Lines2JSONString(&msglist)
}

func Lines2JSONString(multiline *[]string) string {
	j, _ := json.Marshal(multiline)
	res := string(j)
	return res
}