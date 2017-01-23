package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	in := "../examples/hello_world.apib"

	f, err := os.Open(in)

	assert.Nil(t, err, "File failed to load")

	defer f.Close()
	bpParser := Blueprint{}
	_, err = bpParser.Parse(f)

	assert.Nil(t, err, "Blueprint parsing failed")
}
