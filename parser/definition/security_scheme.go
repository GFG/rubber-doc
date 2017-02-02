package definition

// SecurityScheme
type SecurityScheme struct {
	Name        string
	Description string
	Type        string
	Transaction Transaction
	Settings    map[string]interface{}
	Other       map[string]string
}
