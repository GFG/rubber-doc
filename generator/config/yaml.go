package config

import (
	"io/ioutil"

	"github.com/gigforks/yaml"
	"github.com/pkg/errors"
)

// YAML Represents configuration in a yaml file
type YAML struct {
	Combine        bool   `yaml:"combined"`
	SrcDir         string `yaml:"srcDir"`
	DstDir         string `yaml:"dstDir"`
	OutputFilename string `yaml:"output,omitempty"`
	TemplateFiles  []struct {
		SrcFilename string `yaml:"src"`
		DstFilename string `yaml:"dst,omitempty"`
	} `yaml:"templates"`
}

// FromYaml Returns configuration fetched from a yaml file
func FromYaml(filename string) (cfg Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		err = errors.Wrapf(err, "Cannot read from the config file %s", filename)
		return
	}

	y := new(YAML)

	err = yaml.Unmarshal(data, y)
	if err != nil {
		err = errors.Wrapf(err, "Cannot parse the config file %s", filename)
		return
	}

	cfg = NewConfig(y.Combine, y.SrcDir, y.DstDir, y.OutputFilename, y.templates())

	return
}

// templates Returns the configuration for templates
func (y YAML) templates() (config []TemplateConfig) {
	for _, tmpl := range y.TemplateFiles {
		config = append(config, NewTemplateConfig(tmpl.SrcFilename, tmpl.DstFilename))
	}
	return
}
