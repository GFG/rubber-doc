package config

import (
	"io/ioutil"

	"path/filepath"

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

	if err = y.applyAbsToDirectories(filename); err != nil {
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

// applyAbsToDirectories Applies the absolute path of the config's file to source/destination directories
func (y *YAML) applyAbsToDirectories(filename string) (err error) {
	var abs string

	// Fetch absolute path of the configuration file
	if abs, err = filepath.Abs(filepath.Dir(filename)); err != nil {
		return
	}

	// Adds the absolute path to the source directory
	y.SrcDir = filepath.Join(abs, y.SrcDir)
	// Adds the absolute path to the destination directory
	y.DstDir = filepath.Join(abs, y.DstDir)

	return
}
