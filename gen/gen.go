package gen

import (
	"encoding/json"
	"fmt"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

func LoadSpec(reader io.Reader, t SpecType) (spec *Spec, err error) {
	spec = &Spec{}
	switch t {
	case SpecTypeJSON:
		d := json.NewDecoder(reader)
		err = d.Decode(spec)
	case SpecTypeYAML:
		d := yaml.NewDecoder(reader)
		err = d.Decode(spec)
	}
	spec.specType = t
	return
}

func LoadSwagger(spec *Spec) (s *model.Swagger, err error) {
	if spec == nil {
		return
	}
	s = &model.Swagger{
		Operations: make([]model.Operation, 0),
	}
	s.Description = spec.Info.Description
	s.Title = spec.Info.Title
	for path, pathInfoMap := range spec.Paths {
		for method, pathInfo := range pathInfoMap {
			comment := pathInfo.Summary
			if len(comment) == 0 {
				comment = pathInfo.Description
			}
			op := model.Operation{
				Comment:     comment,
				OperationID: pathInfo.OperationId,
				APIMethodUC: strings.ToUpper(method),
				APIMethodLC: strings.ToLower(method),
				APIPath:     path,
				Parameters:  make([]model.Parameter, 0),
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
				p := model.Parameter{
					Name:    specParam.Name,
					Names:   []model.ParameterName{{Swagger: specParam.Name, JS: strcase.ToLowerCamel(specParam.Name)}},
					Comment: specParam.Desc,
					Type:    model.ParameterType(specParam.In),
				}
				if p.Names[0].JS != p.Name {
					p.Names = append(p.Names, model.ParameterName{JS: p.Name, Swagger: p.Name})
				}
				switch p.Type {
				case model.ParameterTypeHeader:
					p.TypeIs.Header = true
				case model.ParameterTypePath:
					p.TypeIs.Path = true
				case model.ParameterTypeQuery:
					p.TypeIs.Query = true
				case model.ParameterTypeBody:
					p.TypeIs.Body = true
				case model.ParameterTypeForm:
					p.TypeIs.FormData = true
				}
				op.Parameters = append(op.Parameters, p)
			}
			s.Operations = append(s.Operations, op)
		}
	}
	return
}
