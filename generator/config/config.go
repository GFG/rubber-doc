package config

import "path/filepath"

// Config Represents the configuration used on the generation process
type config struct {
	combine        bool
	srcDir         string
	dstDir         string
	outputFilename string
	templates      []TemplateConfig
}

// NewConfig Return an instance of configuration
func NewConfig(combine bool, srcDir string, dstDir string, outputFilename string, tmpls []TemplateConfig) (cfg config) {
	return config{
		combine,
		srcDir,
		dstDir,
		outputFilename,
		tmpls,
	}
}

// IsCombined Returns the value for the flag combined
func (c config) IsCombined() bool {
	return c.combine
}

// Src Returns the directory where the templates are located
func (c config) Src() string {
	return c.srcDir
}

// Dst Returns the output's directory
func (c config) Dst() string {
	return c.dstDir
}

// Output Returns the file's absolute path containing the combined output
func (c config) Output() string {
	return filepath.Join(c.dstDir, c.outputFilename)
}

// Templates Returns the configuration for templates
func (c config) Templates() []TemplateConfig {
	return c.templates
}
