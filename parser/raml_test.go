package parser

import (
	"testing"

	"os"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
	"github.com/stretchr/testify/assert"
)

var apiDef *definition.Api

func TestMain(m *testing.M) {
	var err error

	p := NewRamlParser()

	apiDef, err = p.Parse("testdata/raml/api.raml", transformer.NewRamlTransformer())

	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestRamlParser_Integration(t *testing.T) {
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
			Name:        "Example",
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

	var min = float64(0.0)

	expectedTraits := []definition.Trait{
		{
			Name:  "paged",
			Usage: "Applies limit and offset parameters to pagination purposes",
			Href: definition.Href{
				Path: "",
				Parameters: []definition.Parameter{
					{
						Name:        "limit",
						Description: "Description for limit property",
						Type:        "integer",
						Required:    false,
						Pattern:     (*string)(nil),
						MinLength:   (*int)(nil),
						MaxLength:   (*int)(nil),
						Min:         &min,
						Max:         (*float64)(nil),
						Example:     nil,
					},
				},
			},
			Transactions: []definition.Transaction{
				{
					Request:  definition.Request{},
					Response: definition.Response{},
				},
			},
		},
		{
			Name:  "secured",
			Usage: "Applies to methods needing security - do not forget to also add securedBy!",
			Transactions: []definition.Transaction{
				{
					Request: definition.Request{
						Title:       "",
						Description: "",
						Body:        nil,
						Headers: []definition.Header{
							{
								Name:    "Authorization",
								Example: "Bearer czZCaGRSa3F0MzpnWDFmQmF0M2JW",
							},
						},
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
						Body:        nil,
						Headers: []definition.Header{
							{
								Name:        "Authorization",
								Description: "Used to send a valid OAuth 2 access token.",
							},
						},
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
			Name:         "OAuth 1.0",
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

	var min = float64(0.0)

	expectedResourcesGroups := []definition.ResourceGroup{
		{
			Title:       "",
			Description: "",
			Resources: []definition.Resource{
				{
					Href: definition.Href{
						FullPath:   "/v1",
						Path:       "/v1",
						Parameters: nil,
					},
					Title:       "1 - Example",
					Description: "",
					Actions:     nil,
					Resources: []definition.Resource{
						{
							Href: definition.Href{
								FullPath:   "/v1/first",
								Path:       "/first",
								Parameters: nil,
							},
							Title:       "1 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										FullPath:   "/v1/first/example",
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "secured",
										},
										{
											Name: "paged",
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
											Method:      "GET",
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
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
													},
													{
														Name:        "limit",
														Description: "Description for limit property",
														Type:        "integer",
														Required:    false,
														Min:         &min,
														Example:     nil,
													},
												},
											},
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "Example",
																CustomType:  nil,
																MediaType:   "application/json",
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
								FullPath:   "/v1/second",
								Path:       "/second",
								Parameters: nil,
							},
							Title:       "2 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										FullPath:   "/v1/second/example",
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "secured",
										},
									},
									SecuredBy: []definition.Option{
										{
											Name:       "oauth_1_0",
											Parameters: nil,
										},
									},
									Actions: []definition.ResourceAction{
										{
											Title:       "",
											Description: "Request Example",
											Method:      "GET",
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
													},
												},
											},
											SecuredBy: []definition.Option{
												{
													Name:       "oauth_1_0",
													Parameters: nil,
												},
											},
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "Example",
																CustomType:  nil,
																MediaType:   "application/json",
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
						FullPath:   "/v2",
						Path:       "/v2",
						Parameters: nil,
					},
					Title:       "2 - Example",
					Description: "",
					Actions:     nil,
					Resources: []definition.Resource{
						{
							Href: definition.Href{
								FullPath:   "/v2/first",
								Path:       "/first",
								Parameters: nil,
							},
							Title:       "1 - Example",
							Description: "",
							Actions:     nil,
							Resources: []definition.Resource{
								{
									Href: definition.Href{
										FullPath:   "/v2/first/example",
										Path:       "/example",
										Parameters: nil,
									},
									Title:       "",
									Description: "",
									Is: []definition.Option{
										{
											Name: "secured",
										},
										{
											Name: "paged",
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
											Method:      "GET",
											Href: definition.Href{
												Path: "",
												Parameters: []definition.Parameter{
													{
														Name:        "example",
														Description: "Query Parameter Example",
														Type:        "date-only",
														Required:    false,
													},
													{
														Name:        "limit",
														Description: "Description for limit property",
														Type:        "integer",
														Required:    false,
														Min:         &min,
													},
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
											Transactions: []definition.Transaction{
												{
													Request: definition.Request{
														Title:       "",
														Description: "",
														Body:        nil,
														Headers: []definition.Header{
															{
																Name:        "Authorization",
																Description: "",
																Example:     nil,
															},
														},
													},
													Response: definition.Response{
														StatusCode:  200,
														Description: "",
														Headers:     nil,
														Body: []definition.Body{
															{
																Description: "",
																Type:        "Example",
																CustomType:  nil,
																MediaType:   "application/json",
																Example:     "",
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

	assert.Exactly(t, expectedResourcesGroups, apiDef.ResourceGroups)
}
