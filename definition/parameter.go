package definition

// Parameter
type Parameter struct {
	Name        string
	Description string
	Type        string
	Required    bool
	Pattern     *string
	MinLength   *int
	MaxLength   *int
	Min         *float64
	Max         *float64
	Example     interface{}
}
