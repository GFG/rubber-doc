package parser

type Specification struct {
	data interface{}
}

type Formatter interface {
	Format(data interface{}) Specification
}

type APIParser interface {
	Parse(filename string, formatter *Formatter) (spec Specification, err error)
}
