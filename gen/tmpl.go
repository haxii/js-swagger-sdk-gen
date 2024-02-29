package gen

import (
	_ "embed"
	"fmt"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"io"
	"text/template"
)

//go:embed tmpl/index.js.gotmpl
var indexTmplSrc string

var (
	indexTmpl, pkgTmpl *template.Template
)

func init() {
	var err error
	indexTmpl, err = template.New("index").Parse(indexTmplSrc)
	if err != nil {
		panic(fmt.Errorf("fail to make package index.js template with error %s", err))
	}
}

func MakeIndex(swag *model.Swagger, w io.Writer) error {
	return indexTmpl.Execute(w, swag)
}
