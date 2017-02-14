package parser

import (
	"sort"
	"testing"

	"github.com/bradfitz/slice"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
	"github.com/stretchr/testify/assert"
)

var apiDef *definition.Api

func TestRamlParser_Integration(t *testing.T) {
	var err error

	p := NewRamlParser()

	apiDef, err = p.Parse("testdata/raml/api.raml", transformer.NewRamlTransformer())

	assert.Nil(t, err, "Raml parsing failed")
	assert.IsType(t, &definition.Api{}, apiDef)
	assert.NotNil(t, apiDef, "The api definition given is empty")

	t.Run("Title", assertTitle)
	t.Run("Version", assertVersion)
	t.Run("BaseURI", assertBaseURI)
	t.Run("BaseURIParameters", assertBaseURIParameters)
	t.Run("Protocols", assertProtocols)
	t.Run("MediaTypes", assertMediaTypes)
	t.Run("CustomTypes", assertCustomTypes)
	t.Run("Traits", assertTraits)
	t.Run("SecuritySchemes", assertSecuritySchemes)
	t.Run("SecuredBy", assertSecuredBy)
	t.Run("Resources", assertResources)
}

func assertTitle(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "Custom API Title", apiDef.Title)
}

func assertVersion(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "1.0", apiDef.Version)
}

func assertBaseURI(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "//{apiEntryPoint}/rest", apiDef.BaseURI)
}

func assertBaseURIParameters(t *testing.T) {
	t.Parallel()

	expectedBaseURIParameters := []definition.Parameter{
		{
			Name:        "apiEntryPoint",
			Description: "The URL used to consume API",
			Type:        "string",
			Required:    true,
		},
	}

	assert.Exactly(t, expectedBaseURIParameters, apiDef.BaseURIParameters)
}

func assertProtocols(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, []definition.Protocol{definition.Protocol("HTTPS")}, apiDef.Protocols)
}

func assertMediaTypes(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, []definition.MediaType{"application/json"}, apiDef.MediaTypes)
}

func assertCustomTypes(t *testing.T) {
	t.Parallel()

	expectedCustomTypes := []definition.CustomType{
		{
			Name:        "Custom",
			Description: "A custom type",
			Type:        "string",
			Default:     nil,
			Examples:    nil,
		},
		{
			Name:        "useTypes.Example",
			Description: "An example type loaded as Library",
			Type:        "array",
			Default:     nil,
			Examples:    nil,
		},
	}

	assert.Exactly(t, expectedCustomTypes, apiDef.CustomTypes)
}

func assertTraits(t *testing.T) {
	t.Parallel()

	expectedTraits := []definition.Trait{
		{
			Name:  "useTraits.secured",
			Usage: "Apply to methods needing security - do not forget to also add securedBy!",
			Transactions: []definition.Transaction{
				{
					Request: definition.Request{
						Title:       "",
						Description: "",
						Method:      "",
						Body:        nil,
						Headers: []definition.Header{
							{
								Name:    "Authorization",
								Example: "Bearer czZCaGRSa3F0MzpnWDFmQmF0M2JW",
							},
						},
						ContentType: "",
					},
				},
			},
		},
	}

	assert.Exactly(t, expectedTraits, apiDef.Traits)
}

func assertSecuritySchemes(t *testing.T) {
	t.Parallel()

	expectedSecuritySchemes := []definition.SecurityScheme{
		{
			Name:        "oauth_2_0",
			Description: "Custom API supports OAuth 2.0 for authenticating all requests.\n",
			Type:        "OAuth 2.0",
			Transactions: []definition.Transaction{
				{
					Request: definition.Request{
						Title:       "",
						Description: "",
						Method:      "",
						Body:        nil,
						Headers: []definition.Header{
							{
								Name:        "Authorization",
								Description: "Used to send a valid OAuth 2 access token.",
							},
						},
						ContentType: "",
					},
					Response: definition.Response{
						StatusCode:  401,
						Description: "Bad or expired token.",
						Headers:     nil,
						Body:        nil,
					},
				},
			},
			Settings: []definition.SecuritySchemeSetting{
				{
					Name: "accessTokenUri",
					Data: "https://example.net/oauth/access-token",
				},
				{
					Name: "authorizationGrants",
					Data: []interface{}{
						"urn:ietf:params:oauth:grant-type:saml2-bearer",
					},
				},
				{
					Name: "authorizationUri",
					Data: "https://example.net/oauth/authorize",
				},
				{
					Name: "scopes",
					Data: []interface{}{
						"users.read",
					},
				},
			},
		},
		{
			Name:         "useSecuritySchemes.oauth_1_0",
			Description:  "OAuth 1.0 loaded as Library",
			Type:         "OAuth 1.0",
			Transactions: nil,
			Settings: []definition.SecuritySchemeSetting{
				{
					Name: "authorizationUri",
					Data: "https://example.com/oauth/authorize",
				},
				{
					Name: "requestTokenUri",
					Data: "https://api.example.com/oauth/request_token",
				},
				{
					Name: "tokenCredentialsUri",
					Data: "https://api.example.com/oauth/access_token",
				},
			},
		},
	}

	// Due to transformation of map[string]interface into array
	// we need to sort before comparing structures
	slice.Sort(apiDef.SecuritySchemes, func(i, j int) bool {
		isSorted := []string{
			apiDef.SecuritySchemes[i].Name,
			apiDef.SecuritySchemes[j].Name,
		}

		return sort.StringsAreSorted(isSorted)
	})

	for _, scheme := range apiDef.SecuritySchemes {
		slice.Sort(scheme.Settings, func(i, j int) bool {
			isSorted := []string{
				scheme.Settings[i].Name,
				scheme.Settings[j].Name,
			}

			return sort.StringsAreSorted(isSorted)
		})
	}
	// END

	assert.Exactly(t, expectedSecuritySchemes, apiDef.SecuritySchemes)
}

func assertSecuredBy(t *testing.T) {
	t.Parallel()

	expectedSecuredBy := []definition.Option{
		{
			Name:       "oauth_2_0",
			Parameters: map[string]interface{}{"scopes": []interface{}{}},
		},
	}

	assert.Exactly(t, expectedSecuredBy, apiDef.SecuredBy)
}

func assertResources(t *testing.T) {
	t.Parallel()

	expectedResourcesGroups := []definition.ResourceGroup{
		{
			Title:       "",
			Description: "",
			Resources: []definition.Resource{
				{
					Href: definition.Href{
						Path:       "/v1",
						Parameters: nil,
					},
					Title:       "1 - Example",
					Description: "",
					Actions:     nil,
					Resources: []definition.Resource{
						{
							Href: definition.Href{
								Path:       "/first",
								Parameters: nil,
							},
							Title:       "1 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "useTraits.secured",
										},
									},
									SecuredBy: []definition.Option{
										{
											Name: "oauth_2_0",
											Parameters: map[string]interface{}{
												"scopes": []interface{}{
													"users.read",
												},
											},
										},
									},
									Actions: []definition.ResourceAction{
										{
											Title:       "",
											Description: "Request Example",
											Method:      "",
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
														Pattern:     (*string)(nil),
														MinLength:   (*int)(nil),
														MaxLength:   (*int)(nil),
														Min:         (*float64)(nil),
														Max:         (*float64)(nil),
														Example:     nil,
													},
												},
											},
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Method:      "GET",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
														ContentType: "",
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "",
																CustomType: definition.CustomType{
																	Name:        "",
																	Description: "",
																	Type:        "useTypes.Example",
																	Default:     nil,
																	Properties:  nil,
																	Examples:    nil,
																},
																MediaType: "application/json",
																Example:   "",
															},
														},
													},
												},
											},
										},
									},
									Resources: nil,
								},
							},
						},
						{
							Href: definition.Href{
								Path:       "/second",
								Parameters: nil,
							},
							Title:       "2 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "useTraits.secured",
										},
									},
									SecuredBy: []definition.Option{
										{
											Name:       "useSchemes.oauth_1_0",
											Parameters: nil,
										},
									},
									Actions: []definition.ResourceAction{
										{
											Title:       "",
											Description: "Request Example",
											Method:      "",
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
														Pattern:     (*string)(nil),
														MinLength:   (*int)(nil),
														MaxLength:   (*int)(nil),
														Min:         (*float64)(nil),
														Max:         (*float64)(nil),
														Example:     nil,
													},
												},
											},
											SecuredBy: nil,
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Method:      "GET",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
														ContentType: "",
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "",
																CustomType: definition.CustomType{
																	Name:        "",
																	Description: "",
																	Type:        "useTypes.Example",
																	Default:     nil,
																	Properties:  nil,
																	Examples:    nil,
																},
																MediaType: "application/json",
																Example:   "",
															},
														},
													},
												},
											},
										},
									},
									Resources: nil,
								},
							},
						},
					},
				},
				{
					Href: definition.Href{
						Path:       "/v2",
						Parameters: nil,
					},
					Title:       "2 - Example",
					Description: "",
					Actions:     nil,
					Resources: []definition.Resource{
						{
							Href: definition.Href{
								Path:       "/first",
								Parameters: nil,
							},
							Title:       "1 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "useTraits.secured",
										},
									},
									SecuredBy: []definition.Option{
										{
											Name: "oauth_2_0",
											Parameters: map[string]interface{}{
												"scopes": []interface{}{
													"users.read",
												},
											},
										},
									},
									Actions: []definition.ResourceAction{
										{
											Title:       "",
											Description: "Request Example",
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
														Pattern:     (*string)(nil),
														MinLength:   (*int)(nil),
														MaxLength:   (*int)(nil),
														Min:         (*float64)(nil),
														Max:         (*float64)(nil),
														Example:     nil,
													},
												},
											},
											SecuredBy: nil,
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Method:      "GET",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
														ContentType: "",
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "",
																CustomType: definition.CustomType{
																	Name:        "",
																	Description: "",
																	Type:        "useTypes.Example",
																	Default:     nil,
																	Properties:  nil,
																	Examples:    nil,
																},
																MediaType: "application/json",
																Example:   "",
															},
														},
													},
												},
											},
										},
									},
									Resources: nil,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, groups := range apiDef.ResourceGroups {
		// Due to transformation of map[string]interface into array
		// we need to sort before comparing structures
		slice.Sort(groups.Resources, func(i, j int) bool {
			isSorted := []string{
				groups.Resources[i].Href.Path,
				groups.Resources[j].Href.Path,
			}
			return sort.StringsAreSorted(isSorted)
		})
	}

	assert.Exactly(t, expectedResourcesGroups, apiDef.ResourceGroups)
}
