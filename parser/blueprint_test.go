package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlueprintParser(t *testing.T) {
	p := NewBlueprintParser()
	f := new(BlueprintFormatter)

	spec, err := p.Parse("testdata/blueprint/simple.apib", *f)

	assert.Nil(t, err, "Blueprint parsing failed")
	assert.IsType(t, Specification{}, spec)
}
