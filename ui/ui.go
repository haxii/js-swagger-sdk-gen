package ui

import (
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed index.gohtml
var indexTmplSrc string

var (
	indexTmpl *template.Template
)

func init() {
	var err error
	indexTmpl, err = template.New("index").Parse(indexTmplSrc)
	if err != nil {
		panic(fmt.Errorf("fail to make swagger ui index.html template with error %s", err))
	}
}
