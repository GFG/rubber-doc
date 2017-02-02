package transformer

import (
	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
)

type RamlTransformer struct {
	def *definition.Api
}

func NewRamlTransformer() Transformer {
	return new(RamlTransformer)
}

func (f *RamlTransformer) Transform(data interface{}) (def *definition.Api) {
	d, ok := data.(raml.APIDefinition)
	if !ok {
		return
	}

	f.def = &definition.Api{}

	f.title(d)
	f.version(d)
	f.baseUri(d)

	return f.def
}

func (f *RamlTransformer) title(def raml.APIDefinition) {
	f.def.Title = def.Title
}

func (f *RamlTransformer) version(def raml.APIDefinition) {
	f.def.Version = def.Version
}

func (f *RamlTransformer) baseUri(def raml.APIDefinition) {
	f.def.BaseURI = def.BaseUri
}
