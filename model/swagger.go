package model

import (
	"fmt"
	"net/url"
)

type ParameterType string

const (
	ParameterTypeHeader ParameterType = "header"
	ParameterTypePath   ParameterType = "path"
	ParameterTypeQuery  ParameterType = "query"
	ParameterTypeBody   ParameterType = "body"
	ParameterTypeForm   ParameterType = "formData"
)

type ParameterName struct {
	Swagger, JS string
}

func (p ParameterName) SwaggerVarInPath() string {
	return fmt.Sprintf("{%s}", p.Swagger)
}

type Parameter struct {
	Name    string
	Names   []ParameterName // with normalized name, in camel case
	Comment string
	Type    ParameterType
	TypeIs  struct {
		Header, Path, Query, Body, FormData bool
	}
}

type Operation struct {
	Comment     string
	OperationID string
	APIMethodUC string // method upper case
	APIMethodLC string // method lower case
	APIPath     string
	Parameters  []Parameter
}

type UICdn string

const (
	UICdnJsCdn     UICdn = "https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/"
	UICdnBootCdn   UICdn = "https://cdn.bootcdn.net/ajax/libs/swagger-ui/"
	UICdnWebStatic UICdn = "https://cdnjs.webstatic.cn/ajax/libs/swagger-ui/"
)

type UI struct {
	Version string
	CDN     UICdn
}

func (u *UI) Norm() {
	if len(u.CDN) == 0 {
		u.CDN = UICdnWebStatic
	}
	if len(u.Version) == 0 {
		u.Version = "5.11.8"
	}
}
func (u *UI) CSS() string {
	u.Norm()
	s, _ := url.JoinPath(string(u.CDN), "/", u.Version, "swagger-ui.css")
	return s
}

func (u *UI) BundleJS() string {
	u.Norm()
	s, _ := url.JoinPath(string(u.CDN), "/", u.Version, "swagger-ui-bundle.js")
	return s
}

type FileType int

const (
	FileTypeJSON FileType = iota
	FileTypeYAML
)

func (t FileType) String() string {
	switch t {
	case FileTypeJSON:
		return "swagger.json"
	case FileTypeYAML:
		return "swagger.yaml"
	}
	return ""
}

type Swagger struct {
	JSPackage struct {
		Name     string
		Version  string
		Author   string
		License  string
		CommonJS bool // set to true to output cjs, otherwise ejs
	}

	UI UI

	Description string
	Title       string
	Operations  []Operation

	FileType   FileType
	RawContent []byte
}
