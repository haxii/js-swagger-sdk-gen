package registry

import (
	"fmt"
	"testing"
)

func TestGetNpmRc(t *testing.T) {
	rcs, err := GetNpmRC()
	if err != nil {
		t.Fatal(err)
	}
	for _, rc := range rcs {
		fmt.Println(rc.FilePath)
		fmt.Println(rc.DefaultURL)
		fmt.Println(rc.Tokens)
		fmt.Println("------")
	}

	fmt.Println(GetDefaultNpmRC())
}
