package ui

import (
	"github.com/haxii/js-swagger-sdk-gen/model"
	"testing"
)

func TestMakeUI(t *testing.T) {
	if err := MakeUI(&model.Swagger{
		Info: model.SwaggerInfo{
			Description: `desc with "" in content`,
			Title:       "title",
		},
		FileType:   model.SwaggerFileTypeJSON,
		RawContent: []byte(`{"content":"swagger raw content"}`),
	}, "/tmp/swagger-ui"); err != nil {
		t.Fatal(err)
	}
}
