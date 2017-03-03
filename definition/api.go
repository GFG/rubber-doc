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
	CustomTypes       []CustomType
	Traits            []Trait
	SecuritySchemes   []SecurityScheme
	SecuredBy         []Option
	ResourceGroups    []ResourceGroup
}

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

// Trait Optional definition that
type Trait struct {
	Name         string
	Usage        string
	Description  string
	Protocols    []Protocol
	Href         Href
	Transactions []Transaction
}

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
	Body        []Body
	Headers     []Header
	ContentType string
}

// Response
type Response struct {
	StatusCode  int
	Description string
	Headers     []Header
	Body        []Body
}

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

// Href
type Href struct {
	Path       string
	Parameters []Parameter
}

// Header
type Header struct {
	Name        string
	Description string
	Example     interface{}
}

// Body
type Body struct {
	Description string
	Type        string
	CustomType  CustomType
	MediaType   MediaType
	Example     string
}

// Option
type Option struct {
	Name       string
	Parameters map[string]interface{}
}
