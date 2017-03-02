package generator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate_Integration(t *testing.T) {
	// Clean up
	defer func() {
		os.RemoveAll("/tmp/rubberdoc")
	}()

	checks := []struct {
		ApiDef         definition.Api
		ConfigFile     string
		Output         string
		ExpectedOutput string
	}{
		{
			definition.Api{
				Title:      "Custom API Title",
				Version:    "1.0",
				BaseURI:    "https://api.example.com",
				Protocols:  []definition.Protocol{"HTTPS", "HTTP"},
				MediaTypes: []definition.MediaType{"application/json"},
			},
			"testdata/html/config.yaml",
			"/tmp/rubberdoc/output/html/simple.html",
			"testdata/html/simple.html",
		},
		{
			definition.Api{
				Title:      "Custom API Title",
				Version:    "1.0",
				BaseURI:    "https://api.example.com",
				Protocols:  []definition.Protocol{"HTTPS", "HTTP"},
				MediaTypes: []definition.MediaType{"application/json"},
			},
			"testdata/html/advanced/config.yaml",
			"/tmp/rubberdoc/output/html/advanced/index.html",
			"testdata/html/advanced/index.html",
		},
	}

	for _, check := range checks {
		gen, err := NewGenerator(check.ConfigFile, check.ApiDef)
		assert.Nil(t, err)

		err = gen.Generate()
		assert.Nil(t, err)

		expected, err := testLoadFile(check.ExpectedOutput)
		assert.Nil(t, err)

		output, err := testLoadFile(check.Output)
		assert.Nil(t, err)

		assert.Exactly(t, expected, output)
	}
}

// testLoadFile Reads the content of a file and returns it as string
func testLoadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}
