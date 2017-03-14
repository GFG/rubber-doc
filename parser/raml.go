package parser

import (
	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
)

//RamlParser  Concrete's parser definition
type RamlParser struct{}

//NewRamlParser Creates a raml parser
func NewRamlParser() Parser {
	return &RamlParser{}
}

//Parse Concrete implementation of the Parser.Parse method
func (rp RamlParser) Parse(filename string, tra transformer.Transformer) (def *definition.Api, err error) {
	ramlDef := new(raml.APIDefinition)

	if err = raml.ParseFile(filename, ramlDef); err != nil {
		return
	}

	def, err = tra.Transform(*ramlDef)

	return
}
