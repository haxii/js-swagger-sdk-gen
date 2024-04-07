package model

import (
	"fmt"
	"net/url"
	"path"
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

	Tags         []string
	HasAPIDocURL bool
	APIDocURL    string // doc's URL
}

type SwaggerUICdn string

const (
	SwaggerUICdnJsCdn     SwaggerUICdn = "https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/"
	SwaggerUICdnBootCdn   SwaggerUICdn = "https://cdn.bootcdn.net/ajax/libs/swagger-ui/"
	SwaggerUICdnWebStatic SwaggerUICdn = "https://cdnjs.webstatic.cn/ajax/libs/swagger-ui/"
)

type SwaggerUI struct {
	Version string
	CDN     SwaggerUICdn
}

func (u *SwaggerUI) Norm() {
	if len(u.CDN) == 0 {
		u.CDN = SwaggerUICdnWebStatic
	}
	if len(u.Version) == 0 {
		u.Version = "5.11.8"
	}
}
func (u *SwaggerUI) CSS() string {
	u.Norm()
	s, _ := url.JoinPath(string(u.CDN), "/", u.Version, "swagger-ui.css")
	return s
}

func (u *SwaggerUI) BundleJS() string {
	u.Norm()
	s, _ := url.JoinPath(string(u.CDN), "/", u.Version, "swagger-ui-bundle.js")
	return s
}

type SwaggerFileType int

const (
	SwaggerFileTypeJSON SwaggerFileType = iota
	SwaggerFileTypeYAML
)

func (t SwaggerFileType) String() string {
	switch t {
	case SwaggerFileTypeJSON:
		return "swagger.json"
	case SwaggerFileTypeYAML:
		return "swagger.yaml"
	}
	return ""
}

type SwaggerInfo struct {
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
	Title       string `json:"title" yaml:"title"`
	Contact     struct {
		Name  string `json:"name" yaml:"name"` // ext: contact's name
		Email string `json:"email" yaml:"email"`
	} `json:"contact" yaml:"contact"`
	License struct {
		Name string `json:"name" yaml:"name"`
	} `json:"license" yaml:"license"`
	Homepage    string `json:"x-homepage" yaml:"x-homepage"`       // ext: swagger UI homepage
	PackageName string `json:"x-package-name" yaml:"package-name"` // ext: js package name
}

type SwaggerGenConf struct {
	CommonJS        bool // set to true to output cjs, otherwise ejs
	UrlRefInComment bool // set swagger UI deep link url in js comment, like https://petstore.swagger.io/#/pet/getPetById
}

type Swagger struct {
	GenConf SwaggerGenConf // config options

	PkgJSON PackageJSON // lib's package.json

	UI       SwaggerUI
	Info     SwaggerInfo
	FileType SwaggerFileType

	Operations []Operation

	RawContent []byte
}

// SetUrlRefInComment set UrlRefInComment in Swagger.GenConf as true and
func (swag *Swagger) SetUrlRefInComment() error {
	homepage := swag.PkgJSON.Homepage
	if len(homepage) == 0 {
		return nil
	}
	baseURL, err := url.Parse(homepage)
	if err != nil {
		return err
	}
	swag.GenConf.UrlRefInComment = true
	baseURL = baseURL.JoinPath("/")
	for i, op := range swag.Operations {
		swag.Operations[i].HasAPIDocURL = true
		tag := "default"
		if len(op.Tags) > 0 {
			tag = op.Tags[0]
		}
		baseURL.Fragment = path.Join(tag, op.OperationID)
		swag.Operations[i].APIDocURL = baseURL.String()
	}
	return nil
}
