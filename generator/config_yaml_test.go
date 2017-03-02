package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigYaml(t *testing.T) {
	checks := []struct {
		ConfigFile string
		Expected   *ConfigYaml
	}{
		{
			"testdata/config/simple.yaml",
			&ConfigYaml{
				Combine: false,
				SrcDir:  "source",
				DstDir:  "destination",
				TemplateFiles: []TemplateConfigYaml{
					{
						SrcFilename: "simple.tmpl",
						DstFilename: "simple.html",
					},
				},
			},
		},
		{
			"testdata/config/advanced.yaml",
			&ConfigYaml{
				Combine:        true,
				SrcDir:         "source",
				DstDir:         "destination",
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
		},
	}

	for _, check := range checks {
		config, err := NewConfigYaml(check.ConfigFile)
		assert.Nil(t, err)

		assert.Exactly(t, check.Expected, config)
	}
}
