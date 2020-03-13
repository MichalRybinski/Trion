package schemas

const ProjectJSchemaVersion int = 1
const ProjectJSchema string =`{
    "$schema": "http://json-schema.org/draft-07/schema",
    "$id": "http://example.com/example.json",
    "type": "object",
    "title": "The Project Schema",
    "description": "The basic Project JSON Schema",
    "propertyNames": {
    	"enum": ["_id", "name","type","schema_rev"]
  	},
    "properties": {
        "_id": {
            "$id": "#/properties/_id",
            "type": "string",
            "title": "The Id Schema",
            "description": "An UUID for single Project instance",
            "default": "",
            "examples": [
                "00b46e01-3994-4ac2-939e-2d5052a65961"
            ]
        },
        "name": {
            "$id": "#/properties/name",
            "type": "string",
            "title": "The Name Schema",
            "description": "Project human readable name, used as a route in API, no special characters except \"-\"",
            "default": "",
            "examples": [
                "some-name",
                "project"
            ],
            "pattern": "^[a-z0-9]+(?:-[a-z0-9]+)*$"
        },
        "type": {
            "$id": "#/properties/type",
            "type": "string",
            "title": "The Type Schema",
            "description": "An explanation about the purpose of this instance.",
            "default": "",
            "examples": [
                "corporate"
            ]
        },
        "schema_rev": {
            "$id": "#/properties/schema_rev",
            "type": "integer",
            "title": "The project schema revision/version for last update timestamp",
            "default": 1
        }
    },
    "anyOf": [
      {
        "required" : ["_id","name"]
      },
      {
        "required" : ["type","name"]
      }
    ]
}`

const ProjectFilterJSchema = `{
    "$schema": "http://json-schema.org/draft-07/schema",
    "$id": "http://example.com/example.json",
    "type": "object",
    "title": "The Project Filter Schema",
    "description": "The basic Project Filter JSON Schema",
    "propertyNames": {
    	"enum": ["filter"]
    },
    "properties": {
        "filter": {
            "$id": "#/properties/filter",
            "type": "object",
            "title": "The Filter Schema",
            "description": "Filter for filterin projects",
            "default": "",
            "examples": [
                {
                    "type" : {
                        "$in" : [ "customer", "corporate"]
                    }
                }
            ]
        }
    },
    "required" : [ "filter" ]
}`