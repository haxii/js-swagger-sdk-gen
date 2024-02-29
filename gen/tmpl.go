package gen

import (
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed tmpl/index.js.gotmpl
var indexTmplSrc string

//go:embed tmpl/package.json.gotmpl
var pkgTmplSrc string

var (
	indexTmpl, pkgTmpl *template.Template
)

func init() {
	var err error
	indexTmpl, err = template.New("index").Parse(indexTmplSrc)
	if err != nil {
		panic(fmt.Errorf("fail to make package index.js template with error %s", err))
	}
	pkgTmpl, err = template.New("pkg").Parse(pkgTmplSrc)
	if err != nil {
		panic(fmt.Errorf("fail to make package.json template with error %s", err))
	}
}
