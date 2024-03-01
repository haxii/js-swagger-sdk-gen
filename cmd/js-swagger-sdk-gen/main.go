package main

import (
	"github.com/haxii/js-swagger-sdk-gen/registry"
	"os"
)

var (
	// Build of git, got by LDFLAGS on build
	Build = "dev"
	// Version of git, got by LDFLAGS on build
	Version = "dev"
)

func main() {
	var err error
	if _, err = parseOpt(); err != nil {
		log("fail to parse command options with error %s", err)
		os.Exit(1)
	}
	// validate publish options first
	if opt.Publish {
		fillRegistryToken()
	}
}

func fillRegistryToken() {
	if len(opt.RegistryOptions.URL) == 0 {
		debug("registry url not provided, try to load from .npmrc")
		defaultRC, regErr := registry.GetDefaultNpmRC()
		if regErr != nil {
			log("fail to parse .npmrc with error %s", regErr)
			os.Exit(1)
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
			log("fail to parse .npmrc with error %s", regErr)
			os.Exit(1)
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
