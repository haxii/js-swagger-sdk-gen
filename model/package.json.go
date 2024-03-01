package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/mod/semver"
	"net/url"
	"strings"
)

type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"author"`
	Homepage string `json:"homepage,omitempty"`
	License  string `json:"license,omitempty"`
}

type PackageJSON struct {
	PackageInfo
	Main     string          `json:"main"`
	Module   string          `json:"module"`
	Scripts  json.RawMessage `json:"scripts"`
	Keywords []string        `json:"keywords"`
}

func (p *PackageJSON) FromSwagger(swag *Swagger) error {
	p.PackageInfo = swag.JSPackage
	// ensure package name
	if len(p.Name) == 0 {
		p.Name = swag.Info.PackageName
	}
	if len(p.Name) == 0 {
		return errors.New("empty js package name")
	}
	// ensure package version
	if len(p.Version) == 0 {
		p.Version = swag.Info.Version
	}
	p.Version = strings.TrimLeftFunc(p.Version, func(r rune) bool { // remove any leading v
		return string(r) == "v" || string(r) == "V"
	})
	if !semver.IsValid("v" + p.Version) { // use go mod's semver check, which requires a leading v
		return fmt.Errorf("invalid package version %s", p.Version)
	}
	// other keys
	if len(p.Description) == 0 {
		p.Description = swag.Info.Description
	}
	if len(p.Author.Name) == 0 {
		p.Author.Name = swag.Info.Contact.Name
	}
	if len(p.Author.Email) == 0 {
		p.Author.Email = swag.Info.Contact.Email
	}
	if len(p.Homepage) == 0 {
		p.Homepage = swag.Info.Homepage
	}
	if len(p.License) == 0 {
		p.License = swag.Info.License.Name
	}
	p.Main = "index.js"
	p.Module = "index.m.js"
	p.Scripts = []byte(`{}`)
	p.Keywords = []string{"js-swagger-sdk-gen", "axios"}
	return nil
}

func (p *PackageJSON) ID() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Version)
}

func (p *PackageJSON) TarName() string {
	return fmt.Sprintf("%s.tgz", p.ID())
}

func (p *PackageJSON) URL(base *url.URL) string {
	return base.JoinPath(p.Name, "-", p.TarName()).String()
}
