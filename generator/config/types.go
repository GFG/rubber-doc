package config

// Config
type Config interface {
	IsCombined() bool
	Src() string
	Dst() string
	Output() string
	Templates() []TemplateConfig
}

// TemplateConfig
type TemplateConfig interface {
	Src() string
	Dst() string
}
