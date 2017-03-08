package definition

// CustomType represents custom types
type CustomType struct {
	Name        string
	Description string
	Type        interface{}
	Default     interface{}
	Enum        interface{}
	Properties  map[string]interface{}
	Examples    []interface{}
}
