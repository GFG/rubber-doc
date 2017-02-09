package definition

// MediaType represents the media types available on the API. e.g application/json
type MediaType string

// Api Definition structure
type Api struct {
	Title             string
	Version           string
	BaseURI           string
	BaseURIParameters []Parameter
	Protocols         []Protocol
	MediaTypes        []MediaType
	Types             []Type
	SecuritySchemes   []SecurityScheme
	SecuredBy         map[string]interface{}
	ResourceGroups    []ResourceGroup
}

// Type represents custom types
type Type interface{}

// Transaction groups a pair request/response
type Transaction struct {
	Request  Request
	Response Response
}

// Request
type Request struct {
	Title       string
	Description string
	Method      string
	Body        Body
	Headers     []Header
	ContentType string
}

// Response
type Response struct {
	StatusCode  int
	Description string
	Headers     []Header
	Body        Body
}

// Parameter
type Parameter struct {
	Required    bool
	Description string
	Name        string
	Type        Type
}

// Href
type Href struct {
	Path       string
	Parameters []Parameter
}

// Header
type Header struct {
	Key   string
	Value string
}

// Body
type Body struct {
	ContentType string
	Content     string
}