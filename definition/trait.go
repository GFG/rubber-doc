package definition

// Trait Optional definition that
type Trait struct {
	Name         string
	Usage        string
	Description  string
	Protocols    []Protocol
	Href         Href
	Transactions []Transaction
}
