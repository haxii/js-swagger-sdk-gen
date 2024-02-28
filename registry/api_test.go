package registry

import (
	"net/http"
	"testing"
)

func TestMakeIntegrity(t *testing.T) {
	downloadURL := "https://registry.npmjs.org/swagger-petstore-3-sdk/-/swagger-petstore-3-sdk-1.0.2.tgz"
	resp, err := http.Get(downloadURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	integrity, err := MakeIntegrity(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if integrity.Integrity != "sha512-BbtMVIeKbbvZyG2hKHydkEnfYF5GU/gkPICAcdVdjZnxudChXzwKQ7x8NP6ox1mLZ/TABhtflwOG/DQL5GqPIg==" {
		t.Fatal("invalid Integrity")
	}
	if integrity.SHASum != "ef0559393baacfd0951c184fbdbaaaab0abaf0e9" {
		t.Fatal("invalid shasum")
	}
	t.Log(len(integrity.Base64), integrity.Length)
}

func TestPublish(t *testing.T) {
	downloadURL := "https://registry.npmjs.org/swagger-petstore-3-sdk/-/swagger-petstore-3-sdk-1.0.1.tgz"
	resp, err := http.Get(downloadURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	api, err := NewAPI("https://npm.some-registry.com", "some-token")
	if err != nil {
		t.Fatal(err)
	}
	err = api.Publish(resp.Body, []byte(`{
  "version": "1.0.1",
  "license": "MIT",
  "sideEffects": false,
  "main": "dist/cjs/index.js",
  "typings": "dist/types/index.d.ts",
  "module": "dist/esm/index.js",
  "files": [
    "dist",
    "src"
  ],
  "name": "swagger-petstore-3-sdk",
  "description": "This is a sample Pet Store Server based on the OpenAPI 3.0 specification."
}`))
	if err != nil {
		t.Fatal(err)
	}
}
