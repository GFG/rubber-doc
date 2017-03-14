package definition

import "strings"

// CustomType represents custom types
type CustomType struct {
	Name        string
	Description string
	Type        interface{}
	Default     interface{}
	Enum        interface{}
	Properties  []CustomTypeProperty
	Examples    []interface{}
}

// CustomTypeProperty Represents a property of a custom type
type CustomTypeProperty struct {
	Name        string
	Type        string
	Required    bool
	Description string
	Example     string
	Properties  []CustomTypeProperty
}

// CleanCustomTypeName It responsible for removing expressions
func CleanCustomTypeName(name string) string {
	return strings.Trim(name, "[]")
}
