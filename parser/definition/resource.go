package definition

// ResourceAction represents an action for a resource. e.g. GET /examples or POST /examples
type ResourceAction struct {
	Title        string
	Description  string
	Method       string
	URL          string
	Href         Href
	Transactions []Transaction
}

// Resource represents a resource. e.g. /examples
type Resource struct {
	Title           string
	Description     string
	ResourceActions []*ResourceAction
	Href            Href
}

// ResourceGroup groups logically bound resources
type ResourceGroup struct {
	Title       string
	Description string
	Resources   []*Resource
}
