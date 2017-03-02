package generator

import (
	"io/ioutil"

	"github.com/gigforks/yaml"
	"github.com/pkg/errors"
)

// TemplateConfigYaml Represent the config's definition for templates in YAML
type TemplateConfigYaml struct {
	SrcFilename string `yaml:"src"`
	DstFilename string `yaml:"dst,omitempty"`
}

// ConfigYaml Represents the config's definition in YAML
type ConfigYaml struct {
	Combine        bool                 `yaml:"combined"`
	SrcDir         string               `yaml:"srcDir"`
	DstDir         string               `yaml:"dstDir"`
	OutputFilename string               `yaml:"output,omitempty"`
	TemplateFiles  []TemplateConfigYaml `yaml:"templates"`
}

// NewConfigYaml Return Config parsed from yaml
func NewConfigYaml(filename string) (cfg Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		err = errors.Wrapf(err, "Cannot read from the config file %s", filename)
		return
	}

	cfg = new(ConfigYaml)

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		err = errors.Wrapf(err, "Cannot parse the config file %s", filename)
		return
	}

	return
}

// IsCombined Returns the value for the flag combined
func (c ConfigYaml) IsCombined() bool {
	return c.Combine
}

// Src Returns the location of all templates
func (c ConfigYaml) Src() string {
	return c.SrcDir
}

// Dst Returns the location where the output will be located
func (c ConfigYaml) Dst() string {
	return c.DstDir
}

// Output Return the filename with the result of all parsed templates
func (c ConfigYaml) Output() string {
	return c.OutputFilename
}

// Templates Returns all the templates to be processed
func (c ConfigYaml) TemplatesConfig() (tmpls []TemplateConfig) {
	for _, tmpl := range c.TemplateFiles {
		tmpls = append(tmpls, TemplateConfig{tmpl.SrcFilename, tmpl.DstFilename})
	}
	return
}
