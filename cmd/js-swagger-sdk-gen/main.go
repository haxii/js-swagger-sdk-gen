package main

import (
	"github.com/haxii/js-swagger-sdk-gen/gen"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"github.com/haxii/js-swagger-sdk-gen/registry"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// Build of git, got by LDFLAGS on build
	Build = "dev"
	// Version of git, got by LDFLAGS on build
	Version = "dev"
)

func main() {
	if _, err := parseOpt(); err != nil {
		fatal("fail to parse command options with error %s", err)
	}
	// validate publish options first
	if opt.Publish {
		fillRegistryToken()
	}
	// parse swagger
	swag, err := parseSwagger()
	if err != nil {
		fatal("fail to parse swagger with error %s", err)
	}
	swag.JSPackage = opt.PackageInfo
}

func parseSwagger() (*model.Swagger, error) {
	var r io.ReadCloser
	if strings.HasPrefix(opt.AppOptions.File, "http") { // download from internet
		debug("download from url %s", opt.AppOptions.File)
		resp, err := http.Get(opt.AppOptions.File)
		if err != nil {
			return nil, err
		}
		r = resp.Body
	} else {
		debug("open local file %s", opt.AppOptions.File)
		f, err := os.Open(opt.AppOptions.File)
		if err != nil {
			return nil, err
		}
		r = f
	}
	defer r.Close()
	t := model.SwaggerFileTypeJSON
	if !strings.HasSuffix(strings.ToLower(opt.AppOptions.File), "json") {
		t = model.SwaggerFileTypeYAML
	}
	debug("the spec should be %s, try to parse it", t)
	if spec, err := gen.LoadSpec(r, t); err != nil {
		return nil, err
	} else {
		return gen.LoadSwagger(spec)
	}
}

func fillRegistryToken() {
	if len(opt.RegistryOptions.URL) == 0 {
		debug("registry url not provided, try to load from .npmrc")
		defaultRC, regErr := registry.GetDefaultNpmRC()
		if regErr != nil {
			fatal("fail to parse .npmrc with error %s", regErr)
		}
		debug("find .npmrc {%s} with registry {%s}", defaultRC.Path, defaultRC.URL)
		opt.RegistryOptions.URL = defaultRC.URL
		if len(defaultRC.Token) > 0 {
			opt.RegistryOptions.Token = defaultRC.Token
		}
	} else if len(opt.RegistryOptions.Token) == 0 {
		debug("token of %s not provided, try to find from .npmrc", opt.RegistryOptions.URL)
		opt.RegistryOptions.URL = registry.NormRegURL(opt.RegistryOptions.URL)
		rcs, regErr := registry.GetNpmRC()
		if regErr != nil {
			fatal("fail to parse .npmrc with error %s", regErr)
		}
		opt.RegistryOptions.Token = registry.FindToken(opt.RegistryOptions.URL, rcs)
	}
	maskedToken := opt.RegistryOptions.Token
	if len(maskedToken) > 5 {
		maskedToken = maskedToken[:3] + "..." + maskedToken[len(maskedToken)-2:]
	} else if len(maskedToken) == 0 {
		maskedToken = "<empty>"
	}
	log("use registry %s with token %s", opt.RegistryOptions.URL, maskedToken)
}
