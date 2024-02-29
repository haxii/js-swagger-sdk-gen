package gen

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	"time"
)

func LoadSpec(reader io.Reader, t model.FileType) (spec *Spec, err error) {
	spec = &Spec{}
	b := bytes.Buffer{}
	r := io.TeeReader(reader, &b)
	switch t {
	case model.FileTypeJSON:
		d := json.NewDecoder(r)
		err = d.Decode(spec)
	case model.FileTypeYAML:
		d := yaml.NewDecoder(r)
		err = d.Decode(spec)
	}
	spec.FileType = t
	spec.Raw = b.Bytes()
	return
}

func LoadSwagger(spec *Spec) (s *model.Swagger, err error) {
	if spec == nil {
		return
	}
	s = &model.Swagger{
		FileType:   spec.FileType,
		RawContent: spec.Raw,
		Operations: make([]model.Operation, 0),
	}
	s.Description = spec.Info.Description
	s.Title = spec.Info.Title
	for path, pathInfoMap := range spec.Paths {
		for method, pathInfo := range pathInfoMap {
			comment := pathInfo.Summary
			if len(comment) == 0 {
				comment = pathInfo.Description
			}
			op := model.Operation{
				Comment:     comment,
				OperationID: pathInfo.OperationId,
				APIMethodUC: strings.ToUpper(method),
				APIMethodLC: strings.ToLower(method),
				APIPath:     path,
				Parameters:  make([]model.Parameter, 0),
			}
			if len(op.OperationID) == 0 {
				op.OperationID = strcase.ToCamel(fmt.Sprintf("%s %s", op.APIMethodUC, op.APIPath))
			}
			for _, specParam := range pathInfo.Parameters {
				if len(specParam.Ref) > 0 {
					err = fmt.Errorf("unsuppored $ref in parameter of %s %s", op.APIMethodUC, op.APIPath)
					return
				}
				if len(specParam.Name) == 0 {
					err = fmt.Errorf("invalid parameter name in %s %s", op.APIMethodUC, op.APIPath)
					return
				}
				if len(specParam.In) == 0 {
					err = fmt.Errorf("invalid parameter location for %s in %s %s",
						specParam.Name, op.APIMethodUC, op.APIPath)
					return
				}
				p := model.Parameter{
					Name:    specParam.Name,
					Names:   []model.ParameterName{{Swagger: specParam.Name, JS: strcase.ToLowerCamel(specParam.Name)}},
					Comment: specParam.Desc,
					Type:    model.ParameterType(specParam.In),
				}
				if p.Names[0].JS != p.Name {
					p.Names = append(p.Names, model.ParameterName{JS: p.Name, Swagger: p.Name})
				}
				switch p.Type {
				case model.ParameterTypeHeader:
					p.TypeIs.Header = true
				case model.ParameterTypePath:
					p.TypeIs.Path = true
				case model.ParameterTypeQuery:
					p.TypeIs.Query = true
				case model.ParameterTypeBody:
					p.TypeIs.Body = true
				case model.ParameterTypeForm:
					p.TypeIs.FormData = true
				}
				op.Parameters = append(op.Parameters, p)
			}
			s.Operations = append(s.Operations, op)
		}
	}
	return
}

type nopWriter struct {
	size int64
}

func (w *nopWriter) Write(p []byte) (n int, err error) {
	w.size += int64(len(p))
	return len(p), nil
}

func defaultTarHeader(name string, isFolder bool) *tar.Header {
	now := time.Now()
	h := tar.Header{
		Name:       name,
		ModTime:    now,
		AccessTime: now,
		ChangeTime: now,
	}
	if isFolder {
		h.Mode = 0755
		h.Typeflag = tar.TypeDir
	} else {
		h.Mode = 0644
		h.Typeflag = tar.TypeReg
	}
	return &h
}

func writeTarFile(w *tar.Writer, name string, f func(round int, _w io.Writer) error) error {
	nw := &nopWriter{}
	if err := f(1, nw); err != nil {
		return err
	}
	header := defaultTarHeader(name, false)
	header.Size = nw.size
	if err := w.WriteHeader(header); err != nil {
		return err
	}
	return f(2, w)
}

// Generate uses swag to generate a npm tgz file into w,
// if pw is not null, a copy of package.json will write into it for future usage
func Generate(swag *model.Swagger, w, pw io.Writer) error {
	pkg := &model.PackageJSON{}
	pkg.FromSwagger(swag)

	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	if err := tw.WriteHeader(defaultTarHeader("package/", true)); err != nil {
		return err
	}

	// package.json
	if err := writeTarFile(tw, "package/package.json", func(round int, _w io.Writer) error {
		if round == 1 && pw != nil { // write to pw in 1st round
			_w = io.MultiWriter(_w, pw)
		}
		enc := json.NewEncoder(_w)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		return enc.Encode(&pkg)
	}); err != nil {
		return err
	}

	// common js
	swag.JSPackage.CommonJS = true
	if err := writeTarFile(tw, "package/index.js", func(_ int, _w io.Writer) error {
		return MakeIndex(swag, _w)
	}); err != nil {
		return err
	}

	// es module file
	swag.JSPackage.CommonJS = false
	if err := writeTarFile(tw, "package/index.m.js", func(_ int, _w io.Writer) error {
		return MakeIndex(swag, _w)
	}); err != nil {
		return err
	}

	return tw.Flush()
}
