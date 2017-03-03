package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate_Integration(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	// Clean up
	defer func() {
		os.RemoveAll(outputDir)
	}()

	checks := []struct {
		ApiDef         definition.Api
		Config         ConfigYaml
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
			ConfigYaml{
				Combine: false,
				SrcDir:  "testdata/html",
				DstDir:  filepath.Join(outputDir, "simple"),
				TemplateFiles: []TemplateConfigYaml{
					{
						SrcFilename: "simple.tmpl",
						DstFilename: "index.html",
					},
				},
			},
			filepath.Join(outputDir, "simple/index.html"),
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
			ConfigYaml{
				Combine:        true,
				SrcDir:         "testdata/html/advanced",
				DstDir:         filepath.Join(outputDir, "advanced"),
				OutputFilename: "index.html",
				TemplateFiles: []TemplateConfigYaml{
					{
						SrcFilename: "base.tmpl",
					},
					{
						SrcFilename: "title.tmpl",
					},
					{
						SrcFilename: "version.tmpl",
					},
					{
						SrcFilename: "baseUri.tmpl",
					},
					{
						SrcFilename: "protocols.tmpl",
					},
					{
						SrcFilename: "mediaTypes.tmpl",
					},
				},
			},
			filepath.Join(outputDir, "advanced/index.html"),
			"testdata/html/advanced/index.html",
		},
	}

	for _, check := range checks {
		gen := &Generator{check.ApiDef, check.Config}

		err := gen.Generate()
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
