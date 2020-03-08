package repository

import (
	"time"
	//"fmt"
	//"Trion/common"
	"github.com/MichalRybinski/Trion/schemas"
	//"github.com/xeipuuv/gojsonschema"
)

type ProjectDefinition struct {
	Revision int
	JsonSchema string
	CreatedAt, UpdatedAt time.Time
}

func NewProjectDefinition() *ProjectDefinition {
	return &ProjectDefinition{1, schemas.ProjectJSchema, time.Now(), time.Now()}
} 
//interfejs do struct: func (p *ProjectDefinition) methodname (return values)