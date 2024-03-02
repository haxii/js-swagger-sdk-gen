package main

import (
	"errors"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
	"os"
)

type AppOptions struct {
	File    string `short:"f" long:"file" default:"./swagger.json" description:"swagger file, support both local and remote json/yaml files"`
	Target  string `short:"t" long:"target"  default:"./" description:"target dir to generate the SDK"`
	Publish bool   `short:"p" long:"publish"  description:"publish to the registry directly, if set, the tarball will not write to the target dir"`
	UIDir   string `long:"ui"  description:"generate swagger ui to this dir for distribution"`
}

type MiscOptions struct {
	ShowVer bool `long:"version" description:"display application version" json:"-"`
	Verbose bool `long:"verbose" description:"verbose the output"`
}

type RegistryOptions struct {
	URL   string `long:"registry-url" env:"NPM_REGISTRY_URL" description:"npm registry url to publish the SDK, default from .npmrc"`
	Token string `long:"registry-token" env:"NPM_REGISTRY_TOKEN" description:"npm registry token to publish the SDK, default from .npmrc" yaml:"-" json:"-"`
}

type Options struct {
	AppOptions
	RegistryOptions
	MiscOptions
	model.PackageInfo
}

var (
	// opt
	opt *Options
)

func parseOpt() ([]string, error) {
	opt = &Options{}
	parser := flags.NewParser(&opt.AppOptions, flags.Default)
	parser.ShortDescription = "JavaScript Swagger SDK Generator"
	parser.LongDescription = "Generate and publish a JavaScript SDK using axios with given swagger v2 specification."
	if _, err := parser.AddGroup("SDK Package Options", "", &opt.PackageInfo); err != nil {
		return nil, err
	}
	if _, err := parser.AddGroup("NPM Registry Options", "", &opt.RegistryOptions); err != nil {
		return nil, err
	}
	if _, err := parser.AddGroup("Miscellaneous Options", "", &opt.MiscOptions); err != nil {
		return nil, err
	}
	args, err := parser.Parse()
	if err != nil {
		code := 1
		var fe *flags.Error
		if errors.As(err, &fe) {
			if errors.Is(fe.Type, flags.ErrHelp) {
				code = 0
			}
		}
		os.Exit(code)
	}
	if opt.ShowVer {
		logVer()
		os.Exit(0)
	} else if opt.Verbose {
		logVer()
		debug("use following options to generate npm package")
		enc := yaml.NewEncoder(os.Stdout)
		_ = enc.Encode(opt)
	}
	return args, nil
}

func logVer() {
	log("js-swagger-sdk-gen version %s, build %s\n", Version, Build)
}
