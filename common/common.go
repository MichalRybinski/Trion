package common

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"strings"
	"reflect"

	//"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	m "github.com/MichalRybinski/Trion/common/models"
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
	if IsJSON(desc) {
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
	var err error
	switch v:=input.(type) {
		case []byte: {
			if err = json.Unmarshal(v,&converted); err!=nil {
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
		case m.UserDBModel: converted, err = StructToMap(v,"json")
		case m.AuthModel: converted, err = StructToMap(v,"json")
		default: //nothing, empty filter
	}
	fmt.Println("=> ConvertInterfaceToMapStringInterface, converted: ", converted)
	return converted
}

func MapStringInterface2String(mapVar map[string]interface{}) string {
	msglist := make([]string,0)
	for k,v:= range mapVar {
		msglist=append(msglist,fmt.Sprintf("{\"parameter\": \"%s\", \"value\": \"%s\"}",k,v))
	}
	return "[" + strings.Join(msglist,",") + "]"
}

func Lines2JSONString(multiline *[]string) string {
	j, _ := json.Marshal(multiline)
	res := string(j)
	return res
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func SliceMapToJSONString(itemsMap []map[string]interface{} ) string {
	j, _:= json.Marshal(itemsMap)
	res := string(j)
	return res
}

func MapToJSON(itemMap map[string]interface{} ) []byte {
	j, _:= json.Marshal(itemMap)
	return j
}

// 
func ConvertStructToMap(st interface{}) map[string]interface{} {

	reqRules := make(map[string]interface{})

	v := reflect.ValueOf(st)
	t := reflect.TypeOf(st)

	for i := 0; i < v.NumField(); i++ {
		key := strings.ToLower(t.Field(i).Name)
		typ := v.FieldByName(t.Field(i).Name).Kind().String()
		structTag := t.Field(i).Tag.Get("json")
		jsonName := strings.TrimSpace(strings.Split(structTag, ",")[0])
		value := v.FieldByName(t.Field(i).Name)

		// if jsonName is not empty use it for the key
		if jsonName != "" {
			key = jsonName
		}

		if typ == "string" {
			if !(value.String() == "" && strings.Contains(structTag, "omitempty")) {
				fmt.Println(key, value)
				fmt.Println(key, value.String())
				reqRules[key] = value.String()
			}
		} else if typ == "int" {
			reqRules[key] = value.Int()
		} else {
			reqRules[key] = value.Interface()
		}

	}

	return reqRules
}

func StructToMap(in interface{}, tag string) (map[string]interface{}, error){
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
			v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
			return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
			// gets us a StructField
			fi := typ.Field(i)
			if tagv := fi.Tag.Get(tag); tagv != "" {
					// set key of map to value in struct field
					out[tagv] = v.Field(i).Interface()
			}
	}
	return out, nil
}