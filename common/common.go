package common

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
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

type PostItemResponse struct {
	Success bool `json:"success"`
}

