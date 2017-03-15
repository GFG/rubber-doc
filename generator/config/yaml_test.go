package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromYaml(t *testing.T) {
	checks := []struct {
		ConfigFile string
		Expected   config
	}{
		{
			"testdata/simple.yaml",
			NewConfig(
				false,
				"source",
				"destination",
				"",
				[]TemplateConfig{
					NewTemplateConfig("simple.tmpl", "simple.html"),
				},
			),
		},
		{
			"testdata/advanced.yaml",
			NewConfig(
				true,
				"source",
				"destination",
				"index.html",
				[]TemplateConfig{
					NewTemplateConfig("base.tmpl", ""),
					NewTemplateConfig("title.tmpl", ""),
					NewTemplateConfig("version.tmpl", ""),
					NewTemplateConfig("baseUri.tmpl", ""),
					NewTemplateConfig("protocols.tmpl", ""),
					NewTemplateConfig("mediaTypes.tmpl", ""),
				},
			),
		},
	}

	for _, check := range checks {
		c, err := FromYaml(check.ConfigFile)
		assert.Nil(t, err)

		assert.Exactly(t, check.Expected, c)
	}
}
