package gen

import (
	"github.com/haxii/js-swagger-sdk-gen/model"
	"net/http"
	"os"
	"testing"
)

func TestLoadJSONSpec(t *testing.T) {
	b, err := http.Get("https://petstore.swagger.io/v2/swagger.json")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Body.Close()
	spec, err := LoadSpec(b.Body, model.SwaggerFileTypeJSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(spec.Raw))
	swagger, err := LoadSwagger(spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(swagger)
	pkgTmpl.Execute(os.Stdout, swagger)
}

func TestLoadYAMLSpec(t *testing.T) {
	b, err := http.Get("https://petstore.swagger.io/v2/swagger.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Body.Close()
	spec, err := LoadSpec(b.Body, model.SwaggerFileTypeYAML)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spec)
	swagger, err := LoadSwagger(spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(swagger)
	swagger.GenConf.CommonJS = true
	indexTmpl.Execute(os.Stdout, swagger)
}

func TestGenYAMLSpec(t *testing.T) {
	b, err := http.Get("https://petstore.swagger.io/v2/swagger.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Body.Close()
	spec, err := LoadSpec(b.Body, model.SwaggerFileTypeYAML)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spec)
	swagger, err := LoadSwagger(spec)
	if err != nil {
		t.Fatal(err)
	}
	swagger.JSPackage.Name = "swagger-api"
	swagger.JSPackage.Homepage = "https://petstore.swagger.io/"
	//if err = swagger.SetUrlRefInComment(); err != nil {
	//	t.Fatal(err)
	//}

	target, err := os.Create("/tmp/swagger-api.tgz")
	if err != nil {
		t.Fatal(err)
	}
	defer target.Close()
	if err = Generate(swagger, target, nil); err != nil {
		t.Fatal(err)
	}
}

func TestEmbed(t *testing.T) {
	t.Log(indexTmplSrc)
}
