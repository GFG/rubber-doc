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
			"CustomType's Properties",
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
			"Custom coming from a library",
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
