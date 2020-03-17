package common

import (
  "github.com/ilyakaznacheev/cleanenv"
  //"github.com/hjson/hjson-go"
  //"fmt"
  //"os"
  //"io/ioutil"
)

const HJSONConfFile string = "./configs/config.hjson"
const YmlConfFile string = "./configs/config.yml"
const SysDBName string = "TrionSystem"
const DBUsersCollectionName string = "Users"

type AppConfig struct {
	ServerConfig struct {
		PORT string `yaml:"port" env:"TRION_PORT"` //TODO: both env-upd and some timer with checks for these 2 updated?
		HOST string `yaml:"host" env:"TRION_HOST"`
	} `yaml:"server"`
	DBConfig struct {
		DBType string `yaml:"type" env:"TRION_DB_SRV" env-upd`
		MongoConfig struct {
			URL string			`yaml:"URL" env:"TRION_DB_URL" env-upd`
			ProjectsColl string	`yaml:"ProjectsColl" env:"TRION_PROJECTS_COLL" env-upd`
			SchemasColl string	`yaml:"SchemasColl"  env:"TRION_SCHEMAS_COLL" env-upd`
		} `yaml:"mongoConfig"`
	} `yaml:"db"`
}

//global AppConfig var
var ThisAppConfig = NewAppConfig(YmlConfFile)

func NewAppConfig(confFilePath string) *AppConfig {
	var AppC AppConfig
    err := cleanenv.ReadConfig(confFilePath, &AppC)
    if err != nil { panic(err) }
	return &AppC
}
//TODO update method
/*
func readHJSONFile(fPath string) []byte {
	hjsonFile, err := os.Open(fPath)
	if err != nil { panic(err) }
	byteValue, _ := ioutil.ReadAll(hjsonFile)
	defer hjsonFile.Close()
	return byteValue
}
*/

/*func main() {

    // Now let's look at decoding Hjson data into Go
    // values.
    sampleText := []byte(`
    {
        # specify rate in requests/second
        rate: 1000
        array:
        [
            foo
            bar
        ]
    }`)

    // We need to provide a variable where Hjson
    // can put the decoded data.
    var dat map[string]interface{}

    // Decode and a check for errors.
    if err := hjson.Unmarshal(sampleText, &dat); err != nil {
        panic(err)
    }
    fmt.Println(dat)

    // In order to use the values in the decoded map,
    // we'll need to cast them to their appropriate type.

    rate := dat["rate"].(float64)
    fmt.Println(rate)

    array := dat["array"].([]interface{})
    str1 := array[0].(string)
    fmt.Println(str1)


    // To encode to Hjson with default options:
    sampleMap := map[string]int{"apple": 5, "lettuce": 7}
    hjson, _ := hjson.Marshal(sampleMap)
    // this is short for:
    // options := hjson.DefaultOptions()
    // hjson, _ := hjson.MarshalWithOptions(sampleMap, options)
	fmt.Println(string(hjson))
	
}*/