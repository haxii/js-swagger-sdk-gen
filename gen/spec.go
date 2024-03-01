package gen

import "github.com/haxii/js-swagger-sdk-gen/model"

type SpecParameter struct {
	In   string `json:"in" yaml:"in"`
	Name string `json:"name" yaml:"name"`
	Desc string `json:"description" yaml:"description"`
	Ref  string `json:"$ref" yaml:"$ref"` // ref to def, not supported yet
}

type SpecPath map[string]struct {
	Description string          `json:"description" yaml:"description"`
	Summary     string          `json:"summary" yaml:"summary"`
	OperationId string          `json:"operationId" yaml:"operationId"`
	Tags        []string        `json:"tags" yaml:"tags"`
	Parameters  []SpecParameter `json:"parameters" yaml:"parameters"`
} // key is method

type Spec struct {
	FileType model.SwaggerFileType
	Raw      []byte              // swagger source file
	Info     model.SwaggerInfo   `json:"info" yaml:"info"`
	Paths    map[string]SpecPath `json:"paths" yaml:"paths"` // the key is path
}
