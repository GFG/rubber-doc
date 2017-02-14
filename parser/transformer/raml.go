package transformer

import (
	"strconv"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
)

type RamlTransformer struct{}

func NewRamlTransformer() Transformer {
	return new(RamlTransformer)
}

func (tra *RamlTransformer) Transform(data interface{}) (def *definition.Api) {
	ramlDef, ok := data.(raml.APIDefinition)
	if !ok {
		return
	}

	def = new(definition.Api)

	tra.title(ramlDef, def)
	tra.version(ramlDef, def)
	tra.baseURI(ramlDef, def)
	tra.baseURIParameters(ramlDef, def)
	tra.protocols(ramlDef, def)
	tra.mediaType(ramlDef, def)
	tra.customTypes(ramlDef, def)
	tra.securitySchemes(ramlDef, def)
	tra.securedBy(ramlDef, def)
	tra.resourceGroups(ramlDef, def)
	tra.traits(ramlDef.Traits, def)

	return
}

// title Transforms raml's title definition in api's title definition
func (tra *RamlTransformer) title(ramlDef raml.APIDefinition, def *definition.Api) {
	def.Title = ramlDef.Title
}

// version Transforms raml's version definition in api's version definition
func (tra *RamlTransformer) version(ramlDef raml.APIDefinition, def *definition.Api) {
	def.Version = ramlDef.Version
}

// baseURI Transforms raml's baseURI definition in api's baseURI definition
func (tra *RamlTransformer) baseURI(ramlDef raml.APIDefinition, def *definition.Api) {
	def.BaseURI = ramlDef.BaseURI
}

// baseURIParameters Transforms raml's baseURIParameters definition in api's baseURIParameters definition
func (tra *RamlTransformer) baseURIParameters(ramlDef raml.APIDefinition, def *definition.Api) {
	def.BaseURIParameters = tra.handleParameters(ramlDef.BaseURIParameters)
}

// protocols Transforms raml's protocols definition in api's protocols definition
func (tra *RamlTransformer) protocols(ramlDef raml.APIDefinition, def *definition.Api) {
	def.Protocols = tra.handleProtocols(ramlDef.Protocols)
}

// mediaType Transforms raml's mediaType definition in api's mediaType definition
func (tra *RamlTransformer) mediaType(ramlDef raml.APIDefinition, def *definition.Api) {
	// @todo Apply to raml's parser support for multiple mediaType.
	def.MediaTypes = append(def.MediaTypes, definition.MediaType(ramlDef.MediaType))
}

// customTypes Transforms raml's customTypes definition in api's customTypes definition
func (tra *RamlTransformer) customTypes(ramlDef raml.APIDefinition, def *definition.Api) {
	def.CustomTypes = tra.handleTypes(ramlDef.Types)
}

// securitySchemes Transforms raml's securitySchemes definition in api's securitySchemes definition
func (tra *RamlTransformer) securitySchemes(ramlDef raml.APIDefinition, def *definition.Api) {
	def.SecuritySchemes = tra.handleSecuritySchemes(ramlDef.SecuritySchemes)
}

// securedBy Transforms raml's securedBy definition in api's securedBy definition
func (tra *RamlTransformer) securedBy(ramlDef raml.APIDefinition, def *definition.Api) {
	def.SecuredBy = tra.handleOptions(ramlDef.SecuredBy)
}

// resourceGroups Groups all the raml's resources definition
func (tra *RamlTransformer) resourceGroups(ramlDef raml.APIDefinition, def *definition.Api) {
	//ResourceGroups is an aggregator of resources that is being used by the api definition but not yet supported by RAML.
	def.ResourceGroups = []definition.ResourceGroup{
		{
			Resources: tra.handleResources(ramlDef.Resources),
		},
	}
}

// traits Transforms raml's securitySchemes definition in api's traits definition
func (tra *RamlTransformer) traits(ramlTraits map[string]raml.Trait, def *definition.Api) {
	def.Traits = tra.handleTraits(ramlTraits)
}

// handleProtocols Generic method which handles raml's protocol definition.
func (tra *RamlTransformer) handleProtocols(ramlProtos []string) (protos []definition.Protocol) {
	for _, ramlProto := range ramlProtos {
		protos = append(protos, definition.Protocol(ramlProto))
	}

	return
}

// handleOptions Generic method which handles raml's DefinitionChoice definition.
func (tra *RamlTransformer) handleOptions(ramlOpts []raml.DefinitionChoice) (opts []definition.Option) {
	for _, ramlOpt := range ramlOpts {
		opt := new(definition.Option)

		opt.Name = ramlOpt.Name

		if ramlOpt.Parameters != nil {
			opt.Parameters = ramlOpt.Parameters
		}

		opts = append(opts, *opt)
	}

	return
}

// handleParameters Generic method which handles raml's parameter definition.
func (tra *RamlTransformer) handleParameters(ramlParams map[string]raml.NamedParameter) (params []definition.Parameter) {
	for name, ramlParam := range ramlParams {
		param := new(definition.Parameter)

		// It takes the parameter name over the parameter key from raml definition
		if param.Name = name; ramlParam.Name != "" {
			param.Name = ramlParam.Name
		}

		param.Description = ramlParam.Description
		param.Type = ramlParam.Type
		param.Required = ramlParam.Required
		param.Pattern = ramlParam.Pattern
		param.MinLength = ramlParam.MinLength
		param.MaxLength = ramlParam.MaxLength
		param.Min = ramlParam.Minimum
		param.Max = ramlParam.Maximum
		param.Example = ramlParam.Example

		params = append(params, *param)
	}

	return
}

// handleHeaders Generic method which handles raml's headers definition.
func (tra *RamlTransformer) handleHeaders(ramlHeaders map[raml.HTTPHeader]raml.Header) (headers []definition.Header) {
	for name, ramlHead := range ramlHeaders {
		header := new(definition.Header)

		// It takes the parameter name over the parameter key from raml definition
		if header.Name = string(name); ramlHead.Name != "" {
			header.Name = ramlHead.Name
		}

		header.Description = ramlHead.Description
		header.Example = ramlHead.Example

		headers = append(headers, *header)
	}

	return
}

// handleTypes Generic method which handles raml's type definition.
func (tra *RamlTransformer) handleTypes(ramlTypes map[string]raml.Type) (customTypes []definition.CustomType) {
	for name, ramlType := range ramlTypes {
		customType := definition.CustomType{
			Name:        name,
			Description: ramlType.Description,
			Type:        ramlType.Type,
			Enum:        ramlType.Enum,
			Default:     ramlType.Default,
			Examples:    ramlType.Examples,
		}

		customTypes = append(customTypes, customType)
	}

	return
}

// handleTraits Generic method which handles raml's trait definition.
func (tra *RamlTransformer) handleTraits(ramlTraits map[string]raml.Trait) (traits []definition.Trait) {
	for name, ramlTrait := range ramlTraits {
		trait := definition.Trait{
			Name:        name,
			Usage:       ramlTrait.Usage,
			Description: ramlTrait.Description,
			Href: definition.Href{
				Parameters: tra.handleParameters(ramlTrait.QueryParameters),
			},
			Protocols: tra.handleProtocols(ramlTrait.Protocols),
		}

		var req *definition.Request
		if len(ramlTrait.Headers) != 0 || ramlTrait.Bodies.ApplicationJSON != nil || ramlTrait.Bodies.Type != "" {
			req = &definition.Request{
				Headers: tra.handleHeaders(ramlTrait.Headers),
			}

			req.Body = tra.handleBodies(ramlTrait.Bodies)
		}

		trait.Transactions = tra.handleResponses(req, ramlTrait.Responses)

		traits = append(traits, trait)
	}

	return
}

// handleSecuritySchemes Generic method which handles raml's security schemes definition.
func (tra *RamlTransformer) handleSecuritySchemes(ramlSchemes map[string]raml.SecurityScheme) (schemes []definition.SecurityScheme) {
	for ramlSchemeName, ramlScheme := range ramlSchemes {
		scheme := new(definition.SecurityScheme)

		// It takes the parameter name over the parameter key from raml definition
		if scheme.Name = ramlSchemeName; ramlScheme.Name != "" {
			scheme.Name = ramlScheme.Name
		}

		scheme.Type = ramlScheme.Type
		scheme.Description = ramlScheme.Description

		for name, ramlSet := range ramlScheme.Settings {
			scheme.Settings = append(scheme.Settings, definition.SecuritySchemeSetting{
				Name: name,
				Data: ramlSet,
			})
		}

		var req *definition.Request
		if len(ramlScheme.DescribedBy.Headers) != 0 {
			req = &definition.Request{
				Headers: tra.handleHeaders(ramlScheme.DescribedBy.Headers),
			}
		}

		scheme.Transactions = tra.handleResponses(req, ramlScheme.DescribedBy.Responses)

		schemes = append(schemes, *scheme)
	}

	return
}

// handleResources Generic method which handles raml's resources definition.
func (tra *RamlTransformer) handleResources(ramlResources interface{}) (resources []definition.Resource) {
	switch r := ramlResources.(type) {
	case map[string]raml.Resource:
		for _, ramlRes := range r {
			res := tra.handleResource(ramlRes)

			if ramlRes.Nested != nil {
				res.Resources = tra.handleResources(ramlRes.Nested)
			}

			resources = append(resources, res)
		}
	case map[string]*raml.Resource:
		for _, ramlRes := range r {
			res := tra.handleResource(*ramlRes)

			if ramlRes.Nested != nil {
				res.Resources = tra.handleResources(ramlRes.Nested)
			}

			resources = append(resources, res)
		}
	default:
		//@todo Process better the type of errors
		panic("unsupported type")
	}
	return
}

// handleResource Generic method which handles raml's resource definition.
func (tra *RamlTransformer) handleResource(ramlRes raml.Resource) definition.Resource {
	return definition.Resource{
		Title:       ramlRes.DisplayName,
		Description: ramlRes.Description,
		Href: definition.Href{
			Path:       ramlRes.URI,
			Parameters: tra.handleParameters(ramlRes.URIParameters),
		},
		Is:        tra.handleOptions(ramlRes.Is),
		SecuredBy: tra.handleOptions(ramlRes.SecuredBy),
		Actions:   tra.handleResourceMethods(ramlRes.Methods),
	}
}

// handleResource Generic method which handles raml's method definition.
func (tra *RamlTransformer) handleResourceMethods(ramlMethods []*raml.Method) (actions []definition.ResourceAction) {
	for _, ramlMethod := range ramlMethods {
		action := new(definition.ResourceAction)

		action.Title = ramlMethod.DisplayName
		action.Description = ramlMethod.Description
		action.Href = definition.Href{
			Parameters: tra.handleParameters(ramlMethod.QueryParameters),
		}
		action.Is = tra.handleOptions(ramlMethod.Is)
		action.SecuredBy = tra.handleOptions(ramlMethod.SecuredBy)

		req := &definition.Request{
			Method:  ramlMethod.Name,
			Headers: tra.handleHeaders(ramlMethod.Headers),
		}

		req.Body = tra.handleBodies(ramlMethod.Bodies)

		action.Transactions = tra.handleResponses(req, ramlMethod.Responses)

		actions = append(actions, *action)
	}

	return
}

// handleResponses Generic method which handles raml's responses definition.
// Since raml's specification doesn't consider multiple requests, we have to use the same request definition for each created definition.Transaction
func (tra *RamlTransformer) handleResponses(req *definition.Request, ramlResponses map[raml.HTTPCode]raml.Response) (transactions []definition.Transaction) {
	//Responses are not always present
	if req != nil && len(ramlResponses) == 0 {
		tran := definition.Transaction{
			Request: *req,
		}

		transactions = append(transactions, tran)
	}

	for httpCode, ramlResp := range ramlResponses {
		resp := new(definition.Response)

		// It takes the parameter name over the parameter key from raml definition
		var code = string(httpCode)
		if ramlResp.HTTPCode != "" {
			code = string(ramlResp.HTTPCode)
		}

		//@todo Process the error coming from the conversion
		resp.StatusCode, _ = strconv.Atoi(code)

		resp.Description = ramlResp.Description
		resp.Headers = tra.handleHeaders(ramlResp.Headers)
		resp.Body = tra.handleBodies(ramlResp.Bodies)

		tran := definition.Transaction{
			Request:  *req,
			Response: *resp,
		}

		transactions = append(transactions, tran)
	}

	return
}

// handleBodies Generic method which handles raml's bodies definition.
func (tra *RamlTransformer) handleBodies(ramlBodies raml.Bodies) (bodies []definition.Body) {
	if ramlBodies.ApplicationJSON == nil && ramlBodies.Type == "" {
		return
	}

	body := &definition.Body{
		Type:        ramlBodies.Type,
		Description: ramlBodies.Description,
	}

	if ramlBodies.ApplicationJSON != nil {
		body.MediaType = definition.MediaType("application/json")

		if ramlBodies.Example != "" {
			body.Example = ramlBodies.Example
		}

		cType := &definition.CustomType{
			Type: ramlBodies.ApplicationJSON.Type,
		}

		if ramlBodies.ApplicationJSON.Properties != nil {
			cType.Properties = ramlBodies.ApplicationJSON.Properties
		}

		body.CustomType = *cType
	}

	bodies = append(bodies, *body)

	//@todo Add here your code to support other media types for bodies
	return
}
