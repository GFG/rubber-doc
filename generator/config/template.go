package config

// templateConfig Represents the configuration of a template
type templateConfig struct {
	srcFilename string
	dstFilename string
}

// NewTemplateConfig Returns an instance of configuration of the templates
func NewTemplateConfig(src string, dst string) (cfg TemplateConfig) {
	return templateConfig{src, dst}
}

// Src Returns the absolute path of the template's source
func (t templateConfig) Src() string {
	return t.srcFilename
}

// Dst Returns the absolute path of the template's destination
func (t templateConfig) Dst() string {
	return t.dstFilename
}
