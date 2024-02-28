package internal

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
)

type SpecParameter struct {
	In   string `json:"in" yaml:"in"`
	Name string `json:"name" yaml:"name"`
	Desc string `json:"description" yaml:"description"`
	Ref  string `json:"$ref" yaml:"$ref"` // ref to def
}

type SpecPath map[string]struct {
	Description string          `json:"description" yaml:"description"`
	Summary     string          `json:"summary" yaml:"summary"`
	OperationId string          `json:"operationId" yaml:"operationId"`
	Parameters  []SpecParameter `json:"parameters" yaml:"parameters"`
} // key is method

type Spec struct {
	specType SpecType
	Info     struct {
		Description string `json:"description" yaml:"description"`
	} `json:"info" yaml:"info"`
	Paths   map[string]SpecPath        `json:"paths" yaml:"paths"` // the key is path
	DefJSON map[string]json.RawMessage `json:"definitions" yaml:"-"`
	DefYAML map[string]YamlRawMessage  `json:"-" yaml:"definitions"`
}

// YamlRawMessage from https://github.com/go-yaml/yaml/issues/13#issuecomment-1586325414
type YamlRawMessage struct{ *yaml.Node }

func (n *YamlRawMessage) UnmarshalYAML(node *yaml.Node) error {
	n.Node = node
	return nil
}

type SpecType int

const (
	SpecTypeJSON SpecType = iota
	SpecTypeYAML
)

func LoadSpec(reader io.Reader, t SpecType) (spec *Spec, err error) {
	spec = &Spec{}
	switch t {
	case SpecTypeJSON:
		d := json.NewDecoder(reader)
		err = d.Decode(spec)
	case SpecTypeYAML:
		d := yaml.NewDecoder(reader)
		err = d.Decode(spec)
	}
	spec.specType = t
	return
}
