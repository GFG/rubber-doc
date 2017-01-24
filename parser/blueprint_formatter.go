package parser

type BlueprintFormatter struct {}

func (f BlueprintFormatter) Format (data interface{}) (spec Specification)  {
	if s, ok := data.(interface{}); ok {
		spec.data = s
	}
	return
}
