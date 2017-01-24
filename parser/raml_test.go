package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRamlParser(t *testing.T) {
	p := NewRamlParser()
	f := new(RamlFormatter)

	spec, err := p.Parse("testdata/raml/simple.raml", *f)

	assert.Nil(t, err, "Raml parsing failed")
	assert.IsType(t, Specification{}, spec)
}
