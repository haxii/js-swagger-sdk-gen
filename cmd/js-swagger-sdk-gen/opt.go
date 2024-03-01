package main

import (
	"errors"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"github.com/jessevdk/go-flags"
	"os"
)

type AppOptions struct {
	File    string `short:"f" long:"file" default:"./swagger.json" description:"swagger file, support both local and remote json/yaml files"`
	Target  string `short:"t" long:"target"  default:"./" description:"target dir to generate the SDK"`
	Publish bool   `short:"p" long:"publish"  description:"publish to the registry directly, if set, the tarball will not write to the target dir"`
	GenUI   bool   `long:"ui"  description:"generate swagger ui distribution folder to target dir"`
}

type MiscOptions struct {
	Version bool `long:"version" description:"display application version"`
	Verbose bool `long:"verbose" description:"verbose the output"`
}

type RegistryOptions struct {
	URL   string `long:"registry-url" env:"NPM_REGISTRY_URL" description:"npm registry url to publish the SDK, default from .npmrc"`
	Token string `long:"registry-token" env:"NPM_REGISTRY_TOKEN"  description:"npm registry token to publish the SDK, default from .npmrc"`
}

type Options struct {
	AppOptions
	RegistryOptions
	MiscOptions
	model.PackageInfo
}

func parseOpt() (*Options, error) {
	opts := &Options{}
	parser := flags.NewParser(&opts.AppOptions, flags.Default)
	parser.ShortDescription = "JavaScript Swagger SDK Generator"
	parser.LongDescription = "Generate and publish a JavaScript SDK using axios with given swagger v2 specification."
	if _, err := parser.AddGroup("SDK Package Options", "", &opts.PackageInfo); err != nil {
		return nil, err
	}
	if _, err := parser.AddGroup("NPM Registry Options", "", &opts.RegistryOptions); err != nil {
		return nil, err
	}
	if _, err := parser.AddGroup("Miscellaneous Options", "", &opts.MiscOptions); err != nil {
		return nil, err
	}
	if _, err := parser.Parse(); err != nil {
		code := 1
		var fe *flags.Error
		if errors.As(err, &fe) {
			if errors.Is(fe.Type, flags.ErrHelp) {
				code = 0
			}
		}
		os.Exit(code)
	}
	return opts, nil
}
