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

// CustomTypeByName Returns a CustomType struct based on its name
func (def Api) CustomTypeByName(name string) (ct CustomType) {
	for _, customType := range def.CustomTypes {
		if customType.Name == name {
			ct = customType
			break
		}
	}
	return
}
