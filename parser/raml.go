package parser

import (
	"github.com/Jumpscale/go-raml/raml"
)

type RamlAPIParser struct {}

func NewRamlParser() *RamlAPIParser {
	return &RamlAPIParser{}
}

func (rp RamlAPIParser) Parse(filename string, formatter Formatter) (spec Specification, err error) {
	var data *raml.APIDefinition

	if data, err = raml.ParseFile(filename); err != nil {
		return
	}

	spec = formatter.Format(*data)

	return
}
