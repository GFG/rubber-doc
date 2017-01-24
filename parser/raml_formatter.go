package parser

type RamlFormatter struct {}

func (f RamlFormatter) Format (data interface{}) (spec Specification)  {
	if s, ok := data.(interface{}); ok {
		spec.data = s
	}
	return
}
