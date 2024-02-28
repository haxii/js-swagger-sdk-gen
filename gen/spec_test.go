package gen

import (
	"net/http"
	"testing"
)

func TestLoadJSONSpec(t *testing.T) {
	b, err := http.Get("https://petstore.swagger.io/v2/swagger.json")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Body.Close()
	spec, err := LoadSpec(b.Body, SpecTypeJSON)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spec)
	swagger, err := LoadSwagger(spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(swagger)
}

func TestLoadYAMLSpec(t *testing.T) {
	b, err := http.Get("https://petstore.swagger.io/v2/swagger.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Body.Close()
	spec, err := LoadSpec(b.Body, SpecTypeYAML)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(spec)
	swagger, err := LoadSwagger(spec)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(swagger)
}
