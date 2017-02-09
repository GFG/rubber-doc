package parser

import (
	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
)

// RamlParser
type RamlParser struct{}

// NewRamlParser
func NewRamlParser() Parser {
	return &RamlParser{}
}

// Parse
func (rp RamlParser) Parse(filename string, tra transformer.Transformer) (def *definition.Api, err error) {
	data, err := raml.ParseFile(filename)

	if err != nil {
		return
	}

	def = tra.Transform(*data)

	return
}
