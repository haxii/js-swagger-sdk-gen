package internal

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"strings"
)

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
	APIMethodUC string // method upper case
	APIMethodLC string // method lower case
	APIPath     string
	Parameters  []Parameter
}

type Swagger struct {
	Description string
	Operations  []Operation
}

func LoadSwagger(spec *Spec) (s *Swagger, err error) {
	if spec == nil {
		return
	}
	s = &Swagger{
		Operations: make([]Operation, 0),
	}
	s.Description = spec.Info.Description
	for path, pathInfoMap := range spec.Paths {
		for method, pathInfo := range pathInfoMap {
			comment := pathInfo.Summary
			if len(comment) == 0 {
				comment = pathInfo.Description
			}
			op := Operation{
				Comment:     comment,
				OperationID: pathInfo.OperationId,
				APIMethodUC: strings.ToUpper(method),
				APIPath:     path,
				Parameters:  make([]Parameter, 0),
			}
			if len(op.OperationID) == 0 {
				op.OperationID = strcase.ToCamel(fmt.Sprintf("%s %s", op.APIMethodUC, op.APIPath))
			}
			for _, specParam := range pathInfo.Parameters {
				if len(specParam.Ref) > 0 {
					err = fmt.Errorf("unsuppored $ref in parameter of %s %s", op.APIMethodUC, op.APIPath)
					return
				}
				if len(specParam.Name) == 0 {
					err = fmt.Errorf("invalid parameter name in %s %s", op.APIMethodUC, op.APIPath)
					return
				}
				if len(specParam.In) == 0 {
					err = fmt.Errorf("invalid parameter location for %s in %s %s",
						specParam.Name, op.APIMethodUC, op.APIPath)
					return
				}
				op.Parameters = append(op.Parameters, Parameter{
					Name:     specParam.Name,
					NormName: strcase.ToLowerCamel(specParam.Name),
					Comment:  specParam.Desc,
					Type:     ParameterType(specParam.In),
				})
			}
			s.Operations = append(s.Operations, op)
		}
	}
	return
}
