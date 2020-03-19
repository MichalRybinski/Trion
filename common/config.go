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
const UsersDBName string = "Users"
const UsersDBUsersCollection string = "Users"
const UsersDBAuthCollection string = "Auths"

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
    SecretKey string `yaml:"secretkey" env:"TRION_SECRET_KEY" env-upd`
}

//global AppConfig var
var TrionConfig = NewAppConfig(YmlConfFile)

func NewAppConfig(confFilePath string) *AppConfig {
	var AppC AppConfig
    err := cleanenv.ReadConfig(confFilePath, &AppC)
    if err != nil { panic(err) }
	return &AppC
}