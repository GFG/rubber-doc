package generator

// Config Represents the configuration used by the generator
type Config interface {
	IsCombined() bool
	Src() string
	Dst() string
	Output() string
	TemplatesConfig() []TemplateConfig
}
