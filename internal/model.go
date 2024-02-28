package internal

type ParameterType string

const (
	ParameterTypeHeader ParameterType = "header"
	ParameterTypePath   ParameterType = "path"
	ParameterTypeQuery  ParameterType = "query"
	ParameterTypeBody   ParameterType = "body"
	ParameterTypeForm   ParameterType = "formData"
)

type Parameter struct {
	Name     string
	NormName string // normalized name, in camel case
	Comment  string
	Type     ParameterType
}

type Operation struct {
	Comment     string
	OperationID string
	APIMethod   string
	APIPath     string
	Parameters  []Parameter
}

type Swagger struct {
	Description string
	Operations  []Operation
}
