package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type PackageJSON struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Main        string          `json:"main"`
	Module      string          `json:"module"`
	Scripts     json.RawMessage `json:"scripts"`
	Keywords    []string        `json:"keywords"`
	Author      string          `json:"author"`
	License     string          `json:"license"`
}

func (p *PackageJSON) FromSwagger(swag *Swagger) {
	p.Name = swag.JSPackage.Name
	p.Version = swag.JSPackage.Version
	p.Description = swag.Description
	p.Main = "index.js"
	p.Module = "index.m.js"
	p.Scripts = []byte(`{}`)
	p.Keywords = []string{swag.JSPackage.Name, "js-swagger-sdk-gen", "axios"}
	p.Author = swag.JSPackage.Author
	p.License = swag.JSPackage.License
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
