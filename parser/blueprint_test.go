package parser

import (
	"testing"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
	"github.com/stretchr/testify/assert"
)

type BlueprintParserTest struct {
	apiDef *definition.Api
}

func TestBlueprintParser_Integration(t *testing.T) {
	p := NewBlueprintParser()

	def, err := p.Parse("testdata/blueprint/simple.apib", transformer.NewBlueprintTransformer())


	assert.Nil(t, err, "Blueprint parsing failed")
	assert.IsType(t, &definition.Api{}, def)

	parserTest := &BlueprintParserTest {
		apiDef: def,
	}

	t.Run("Title", parserTest.assertTitle)
	t.Run("Version", parserTest.assertVersion)
	t.Run("BaseURI", parserTest.assertBaseURI)
	t.Run("Protocols", parserTest.assertProtocols)
	t.Run("ResourceGroups", parserTest.assertResourceGroups)
	t.Run("Resources", parserTest.assertResources)
}

func (bp *BlueprintParserTest) assertTitle(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "Real World API", bp.apiDef.Title)
}

func (bp *BlueprintParserTest) assertVersion(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "1.0", bp.apiDef.Version)
}

func (bp *BlueprintParserTest) assertBaseURI(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, "https://alpha-api.app.net", bp.apiDef.BaseURI)
}

func (bp *BlueprintParserTest) assertProtocols(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, []definition.Protocol{definition.Protocol("https")}, bp.apiDef.Protocols)
}

func (bp *BlueprintParserTest) assertResourceGroups(t *testing.T) {
	t.Parallel()
	
	assert.Exactly(t, "Posts", bp.apiDef.ResourceGroups[0].Title)
	assert.Exactly(t, "This section groups App.net post resources.", bp.apiDef.ResourceGroups[0].Description)
	assert.IsType(t, []definition.ResourceGroup{}, bp.apiDef.ResourceGroups)
}

func (bp *BlueprintParserTest) assertResources(t *testing.T) {
	t.Parallel()

	assert.Exactly(t, bp.expectedResources(), bp.apiDef.ResourceGroups[0].Resources)
}

func (bp *BlueprintParserTest) expectedResources() []definition.Resource {
	return []definition.Resource{
		{
			Title: "Post",
			Description: "A Post is the other central object utilized by the App.net Stream API. It has\nrich text and annotations which comprise all of the content a users sees in\ntheir feed. Posts are closely tied to the follow graph...",
			Href: definition.Href{
				Path: "/stream/0/posts/{post_id}",
				Parameters: []definition.Parameter{
					{
						Required: true,
						Description: "The id of the Post.",
						Name: "post_id",
						Example: "1",
						Type: "string",
					},
				},
			},
			Actions: []definition.ResourceAction{
				{
					Title: "Retrieve a Post",
					Description: "Returns a specific Post.",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "GET",
							},
							Response: definition.Response{
								StatusCode: 200,
								Headers: []definition.Header{
									{
										Name: "Content-Type",
										Example: "application/json",
									},
								},
								Body: []definition.Body{
									{
										MediaType: "application/json",
										Example:     "{\n    \"data\": {\n        \"id\": \"1\", // note this is a string\n        \"user\": {\n            ...\n        }\n    },\n    \"meta\": {\n        \"code\": 200,\n    }\n}\n",
									},
								},
							},
						},
					},
				},
				{
					Title: "Delete a Post",
					Description: "Delete a Post. The current user must be the same user who created the Post. It\nreturns the deleted Post on success.",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "DELETE",
							},
							Response: definition.Response{
								StatusCode: 204,
							},
						},
					},
				},
			},
		},
		{
			Title: "Posts Collection",
			Description: "A Collection of posts.",
			Href: definition.Href{
				Path: "/stream/0/posts",
			},
			Actions: []definition.ResourceAction{
				{
					Title: "Create a Post",
					Description: "Create a new Post object. Mentions and hashtags will be parsed out of the post\ntext, as will bare URLs...",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "POST",
								Body: []definition.Body{
									{
										MediaType: "application/json",
										Example:     "{\n    \"data\": {\n        \"id\": \"1\", // note this is a string\n        \"user\": {\n            ...\n        }\n    },\n    \"meta\": {\n        \"code\": 200,\n    }\n}\n",
									},
								},
								Headers: []definition.Header{
									{
										Name: "Content-Type",
										Example: "application/json",
									},
								},
							},
							Response: definition.Response{
								StatusCode: 201,
								Headers: []definition.Header{
									{
										Name: "Content-Type",
										Example: "application/json",
									},
								},
								Body: []definition.Body{
									{
										MediaType: "application/json",
										Example:     "{\n    \"data\": {\n        \"id\": \"1\", // note this is a string\n        \"user\": {\n            ...\n        }\n    },\n    \"meta\": {\n        \"code\": 200,\n    }\n}\n",
									},
								},
							},
						},
					},
				},
				{
					Title: "Retrieve all Posts",
					Description: "Retrieves all posts.",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "GET",
							},
							Response: definition.Response{
								StatusCode: 200,
								Headers: []definition.Header{
									{
										Name: "Content-Type",
										Example: "application/json",
									},
								},
								Body: []definition.Body{
									{
										MediaType: "application/json",
										Example:     "{\n    \"data\": [\n        {\n            \"id\": \"1\", // note this is a string\n            ...\n        },\n        {\n            \"id\": \"2\",\n            ...\n        },\n        {\n            \"id\": \"3\",\n            ...\n        },\n    ],\n    \"meta\": {\n        \"code\": 200,\n    }\n}\n",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Title: "Stars",
			Description: "A User’s stars are visible to others, but they are not automatically added to\nyour followers’ streams.",
			Href: definition.Href{
				Path: "/stream/0/posts/{post_id}/star",
				Parameters: []definition.Parameter{
					{
						Required: true,
						Description: "The id of the Post.",
						Name: "post_id",
						Example: "1",
						Type: "string",
					},
				},
			},
			Actions: []definition.ResourceAction{
				{
					Title: "Star a Post",
					Description: "Save a given Post to the current User’s stars. This is just a “save” action,\nnot a sharing action.\n\n*Note: A repost cannot be starred. Please star the parent Post.*",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "POST",
							},
							Response: definition.Response{
								StatusCode: 200,
								Description: "",
								Headers: []definition.Header{
									{
										Name: "Content-Type",
										Example: "application/json",
									},
								},
								Body: []definition.Body{
									{
										MediaType: "application/json",
										Example:   "{\n    \"data\": {\n        \"id\": \"1\", // note this is a string\n        \"user\": {\n            ...\n        }\n    },\n    \"meta\": {\n        \"code\": 200,\n    }\n}\n",
									},
								},
							},
						},
					},
				},
				{
					Title: "Unstar a Post",
					Description: "Remove a Star from a Post.",
					Transactions: []definition.Transaction{
						{
							Request: definition.Request{
								Method: "DELETE",
							},
							Response: definition.Response{
								StatusCode: 204,
							},
						},
					},
				},
			},
		},
	}
}