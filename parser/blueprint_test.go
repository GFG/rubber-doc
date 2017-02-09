package parser

import (
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
	"github.com/stretchr/testify/assert"
)

func TestBlueprintParser_Parse(t *testing.T) {
	expected := &definition.Api{
		Title: "Real World API",
	}

	p := NewBlueprintParser()

	def, err := p.Parse("testdata/blueprint/simple.apib", transformer.NewBlueprintTransformer())

	assert.Nil(t, err, "Blueprint parsing failed")
	assert.IsType(t, &definition.Api{}, def)
	assert.Equal(t, expected, def)
}
