package generator

// Template Represents the configuration of a template
type TemplateConfig struct {
	SrcFilename string
	DstFilename string
}

// Src Returns the source for the current template
func (t TemplateConfig) Src() string {
	return t.SrcFilename
}

// Src Returns the destination for the current template
func (t TemplateConfig) Dst() string {
	return t.DstFilename
}
