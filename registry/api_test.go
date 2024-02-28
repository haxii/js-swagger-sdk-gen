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
}
