package ui

import (
	_ "embed"
	"fmt"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed tmpl/index.gohtml
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

func MakeUI(swagger *model.Swagger, destFolder string) error {
	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return err
	}
	indexPath := filepath.Join(destFolder, "index.html")
	if indexFile, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return err
	} else if err = indexTmpl.Execute(indexFile, swagger); err != nil {
		return err
	}
	swagPath := filepath.Join(destFolder, swagger.FileType.String())
	return os.WriteFile(swagPath, swagger.RawContent, 0644)
}
