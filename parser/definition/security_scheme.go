package definition

// SecurityScheme
type SecurityScheme struct {
	Name         string
	Description  string
	Type         string
	Transactions []Transaction
	Settings     []SecuritySchemeSetting
}

// Setting
type SecuritySchemeSetting struct {
	Name string
	Data interface{}
}
