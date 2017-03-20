package transformer

import (
	"testing"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/stretchr/testify/assert"
)

func TestRamlTransformer_Transform_WrongSpecType(t *testing.T) {
	t.Parallel()

	trans := NewRamlTransformer()

	_, err := trans.Transform(struct{ Title string }{Title: "Testing wrong spec type"})

	assert.NotNil(t, err, "Testing wrong spec type")
}

func TestRamlTransformer_Transform_Title(t *testing.T) {
	t.Parallel()

	spec := raml.APIDefinition{Title: "Title of the sepecification"}
	expected := &definition.Api{Title: "Title of the sepecification"}

	def, err := NewRamlTransformer().Transform(spec)

	if assert.Nil(t, err) {
		assert.Exactly(t, expected.Title, def.Title)
	}
}

func TestRamlTransformer_Transform_Version(t *testing.T) {
	t.Parallel()

	spec := raml.APIDefinition{Version: "v1"}
	expected := &definition.Api{Version: "v1"}

	def, err := NewRamlTransformer().Transform(spec)

	if assert.Nil(t, err) {
		assert.Exactly(t, expected.Version, def.Version)
	}
}

func TestRamlTransformer_Transform_BaseURI(t *testing.T) {
	t.Parallel()

	spec := raml.APIDefinition{BaseURI: "https://api.example.com"}
	expected := &definition.Api{BaseURI: "https://api.example.com"}

	def, err := NewRamlTransformer().Transform(spec)

	if assert.Nil(t, err) {
		assert.Exactly(t, expected.BaseURI, def.BaseURI)
	}
}

func TestRamlTransformer_Transform_CustomTypes(t *testing.T) {
	t.Parallel()

	checks := []struct {
		Name     string
		Spec     raml.APIDefinition
		Expected *definition.Api
	}{
		{
			"Simple Custom Type",
			raml.APIDefinition{
				Types: map[string]raml.Type{
					"Simple": {
						Description: "A simple custom type",
						Type:        "string",
						Default:     "custom_type",
						Example:     "e.g custom_type1",
						Examples: map[string]interface{}{
							"example1": "e.g custom_type2",
						},
					},
				},
			},
			&definition.Api{
				CustomTypes: []definition.CustomType{
					{
						Name:        "Simple",
						Description: "A simple custom type",
						Type:        "string",
						Default:     "custom_type",
						Examples:    []interface{}{"e.g custom_type2", "e.g custom_type1"},
					},
				},
			},
		},
		{
			"CustomType Properties",
			raml.APIDefinition{
				Types: map[string]raml.Type{
					"Example": {
						Description: "An example type",
						Type:        "object",
						Properties: map[string]interface{}{
							"prop1": "string",
							"prop2": map[interface{}]interface{}{
								"type":        "object",
								"required":    false,
								"description": "some desc",
								"properties": map[interface{}]interface{}{
									"prop3": "string",
								},
							},
							"prop4": map[interface{}]interface{}{
								"type":        "array",
								"required":    true,
								"description": "property as array",
							},
						},
					},
				},
			},
			&definition.Api{
				CustomTypes: []definition.CustomType{
					{
						Name:        "Example",
						Description: "An example type",
						Type:        "object",
						Properties: []definition.CustomTypeProperty{
							{
								Name:     "prop1",
								Type:     "string",
								Required: true,
							},
							{
								Name:        "prop2",
								Type:        "object",
								Required:    false,
								Description: "some desc",
								Properties: []definition.CustomTypeProperty{
									{
										Name:     "prop3",
										Type:     "string",
										Required: true,
									},
								},
							},
							{
								Name:        "prop4",
								Type:        "array",
								Required:    true,
								Description: "property as array",
							},
						},
					},
				},
			},
		},
		{
			"Custom from a library",
			raml.APIDefinition{
				Libraries: map[string]*raml.Library{
					"common": {
						Types: map[string]raml.Type{
							"Simple": {
								Type: "string",
							},
						},
					},
				},
			},
			&definition.Api{
				CustomTypes: []definition.CustomType{
					{
						Name: "Simple",
						Type: "string",
					},
				},
			},
		},
		{
			"Sorted custom types and their properties",
			raml.APIDefinition{
				Types: map[string]raml.Type{
					"Example 2": {
						Type: "object",
						Properties: map[string]interface{}{
							"prop4": "string",
							"prop3": "string",
						},
					},
					"Example 1": {
						Type: "object",
						Properties: map[string]interface{}{
							"prop2": "string",
							"prop1": "string",
						},
					},
				},
			},
			&definition.Api{
				CustomTypes: []definition.CustomType{
					{
						Name: "Example 1",
						Type: "object",
						Properties: []definition.CustomTypeProperty{
							{
								Name:     "prop1",
								Type:     "string",
								Required: true,
							},
							{
								Name:     "prop2",
								Type:     "string",
								Required: true,
							},
						},
					},
					{
						Name: "Example 2",
						Type: "object",
						Properties: []definition.CustomTypeProperty{
							{
								Name:     "prop3",
								Type:     "string",
								Required: true,
							},
							{
								Name:     "prop4",
								Type:     "string",
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	for _, check := range checks {
		check := check
		t.Run(check.Name, func(t *testing.T) {
			t.Parallel()

			def, err := NewRamlTransformer().Transform(check.Spec)

			if assert.Nil(t, err) {
				assert.Exactly(t, check.Expected.CustomTypes, def.CustomTypes)
			}
		})
	}
}

func TestRamlTransformer_Transform_Resources(t *testing.T) {
	t.Parallel()

	checks := []struct {
		Name     string
		Spec     raml.APIDefinition
		Expected *definition.Api
	}{
		{
			"Resource FullPath",
			raml.APIDefinition{
				Resources: map[string]raml.Resource{
					"/first": {
						URI: "/first",
						Nested: map[string]*raml.Resource{
							"/second": {
								URI: "/second",
							},
						},
					},
				},
			},
			&definition.Api{
				ResourceGroups: []definition.ResourceGroup{
					{
						Resources: []definition.Resource{
							{
								Href: definition.Href{
									FullPath: "/first",
									Path:     "/first",
								},
								Resources: []definition.Resource{
									{
										Href: definition.Href{
											FullPath: "/first/second",
											Path:     "/second",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"Sorted Resources",
			raml.APIDefinition{
				Resources: map[string]raml.Resource{
					"/a": {
						URI: "/a",
						Nested: map[string]*raml.Resource{
							"/1": {
								URI: "/1",
							},
						},
					},
					"/b": {
						URI: "/b",
						Nested: map[string]*raml.Resource{
							"/1": {
								URI: "/1",
							},
							"/2": {
								URI: "/2",
							},
						},
					},
				},
			},
			&definition.Api{
				ResourceGroups: []definition.ResourceGroup{
					{
						Resources: []definition.Resource{
							{
								Href: definition.Href{
									FullPath: "/a",
									Path:     "/a",
								},
								Resources: []definition.Resource{
									{
										Href: definition.Href{
											FullPath: "/a/1",
											Path:     "/1",
										},
									},
								},
							},
							{
								Href: definition.Href{
									FullPath: "/b",
									Path:     "/b",
								},
								Resources: []definition.Resource{
									{
										Href: definition.Href{
											FullPath: "/b/1",
											Path:     "/1",
										},
									},
									{
										Href: definition.Href{
											FullPath: "/b/2",
											Path:     "/2",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, check := range checks {
		check := check
		t.Run(check.Name, func(t *testing.T) {
			t.Parallel()

			def, err := NewRamlTransformer().Transform(check.Spec)

			if assert.Nil(t, err) {
				assert.Exactly(t, check.Expected.ResourceGroups, def.ResourceGroups)
			}
		})
	}
}
