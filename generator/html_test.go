package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_HTML_Integration(t *testing.T) {
	var (
		err       error
		outputDir string
		gen       Generator
		expected  string
		output    string
	)

	outputDir, err = ioutil.TempDir("", "")
	assert.Nil(t, err)

	// Clean up
	defer func() {
		os.RemoveAll(outputDir)
	}()

	checks := []struct {
		ApiDef         definition.Api
		Config         config.Config
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
			config.NewConfig(
				false,
				"testdata/html",
				filepath.Join(outputDir, "simple"),
				"",
				[]config.TemplateConfig{
					config.NewTemplateConfig("simple.tmpl", "index.html"),
				},
			),
			filepath.Join(outputDir, "simple/index.html"),
			"testdata/html/index.html",
		},
		{
			definition.Api{
				Title:      "Custom API Title",
				Version:    "1.0",
				BaseURI:    "https://api.example.com",
				Protocols:  []definition.Protocol{"HTTPS", "HTTP"},
				MediaTypes: []definition.MediaType{"application/json"},
			}, config.NewConfig(
				true,
				"testdata/html/advanced",
				filepath.Join(outputDir, "advanced"),
				"index.html",
				[]config.TemplateConfig{
					config.NewTemplateConfig("base.tmpl", ""),
					config.NewTemplateConfig("title.tmpl", ""),
					config.NewTemplateConfig("version.tmpl", ""),
					config.NewTemplateConfig("baseUri.tmpl", ""),
					config.NewTemplateConfig("protocols.tmpl", ""),
					config.NewTemplateConfig("mediaTypes.tmpl", ""),
				},
			),
			filepath.Join(outputDir, "advanced/index.html"),
			"testdata/html/advanced/index.html",
		},
	}

	for _, check := range checks {
		gen, err = NewHTMLGenerator(check.Config, check.ApiDef)
		assert.Nil(t, err)

		err = gen.Generate()
		assert.Nil(t, err)

		expected, err = testLoadFile(check.ExpectedOutput)
		assert.Nil(t, err)

		output, err = testLoadFile(check.Output)
		assert.Nil(t, err)

		assert.Exactly(t, expected, output)
	}
}

// testLoadFile Reads the content of a file and returns it as string
func testLoadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	return string(b), err
}
