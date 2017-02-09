package parser

import (
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
	"github.com/stretchr/testify/assert"
)

func TestRamlParser_Parse(t *testing.T) {
	expected := &definition.Api{
		Title:   "Example API",
		Version: "v1",
		BaseURI: "http://localhost/api",
	}

	p := NewRamlParser()

	def, err := p.Parse("testdata/raml/simple.raml", transformer.NewRamlTransformer())

	assert.Nil(t, err, "Raml parsing failed")
	assert.IsType(t, &definition.Api{}, def)
	assert.Equal(t, expected, def)
}
