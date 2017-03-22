package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromYaml(t *testing.T) {
	abs, _ := filepath.Abs("testdata")

	checks := []struct {
		Name       string
		ConfigFile string
		Expected   config
	}{
		{
			"Simple configuration file",
			"testdata/simple.yaml",
			NewConfig(
				false,
				filepath.Join(abs, "source"),
				filepath.Join(abs, "destination"),
				"",
				[]TemplateConfig{
					NewTemplateConfig("simple.tmpl", "simple.html"),
				},
			),
		},
		{
			"Advanced configuration file",
			"testdata/advanced.yaml",
			NewConfig(
				true,
				filepath.Join(abs, "source"),
				filepath.Join(abs, "destination"),
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
		{
			"Configuration file with relative path",
			"testdata/with_relative_path.yaml",
			NewConfig(
				false,
				filepath.Join(abs, "../source"),
				filepath.Join(abs, "../../destination"),
				"",
				[]TemplateConfig{
					NewTemplateConfig("simple.tmpl", "simple.html"),
				},
			),
		},
	}

	for _, check := range checks {
		check := check
		t.Run(check.Name, func(t *testing.T) {
			t.Parallel()

			c, err := FromYaml(check.ConfigFile)

			if assert.Nil(t, err) {
				assert.Exactly(t, check.Expected, c)
			}
		})
	}
}
