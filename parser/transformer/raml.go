package transformer

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Jumpscale/go-raml/raml"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
)

type RamlTransformer struct{}

func NewRamlTransformer() Transformer {
	return new(RamlTransformer)
}

func (tra *RamlTransformer) Transform(data interface{}) (def *definition.Api, err error) {
	ramlDef, ok := data.(raml.APIDefinition)
	if !ok {
		err = errors.New("The data's struct given isn't supported by the RAML's Transformer")
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
	tra.traits(ramlDef.Traits, def)
	tra.libraries(ramlDef.Libraries, def)

	err = tra.resourceGroups(ramlDef, def)

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
	if ramlDef.MediaType != "" {
		def.MediaTypes = append(def.MediaTypes, definition.MediaType(ramlDef.MediaType))
	}
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
func (tra *RamlTransformer) resourceGroups(ramlDef raml.APIDefinition, def *definition.Api) (err error) {
	var resources []definition.Resource
	if resources, err = tra.handleResources(ramlDef.Resources, new(definition.Resource)); err == nil {
		if len(resources) > 0 {
			//ResourceGroups is an aggregator of resources that is being used by the api definition but not yet supported by RAML.
			def.ResourceGroups = append(def.ResourceGroups, definition.ResourceGroup{Resources: resources})
		}
	}
	return
}

// traits Transforms raml's securitySchemes definition in api's traits definition
func (tra *RamlTransformer) traits(ramlTraits map[string]raml.Trait, def *definition.Api) {
	def.Traits = tra.handleTraits(ramlTraits)
}

// libraries Joins some combined (Types and SecuritySchemes) declarations in libraries with the root api definition
func (tra *RamlTransformer) libraries(ramlLibs map[string]*raml.Library, def *definition.Api) {
	for _, ramlLib := range ramlLibs {
		def.CustomTypes = append(def.CustomTypes, tra.handleTypes(ramlLib.Types)...)
		def.SecuritySchemes = append(def.SecuritySchemes, tra.handleSecuritySchemes(ramlLib.SecuritySchemes)...)
		def.Traits = append(def.Traits, tra.handleTraits(ramlLib.Traits)...)

		if ramlLib.Libraries != nil {
			tra.libraries(ramlLib.Libraries, def)
		}
	}

	return
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

		opt.Name = tra.removeLibraryName(ramlOpt.Name)

		if ramlOpt.Parameters != nil {
			opt.Parameters = ramlOpt.Parameters
		}

		opts = append(opts, *opt)
	}

	return
}

// handleParameters Generic method which handles raml's parameter definition.
func (tra *RamlTransformer) handleParameters(ramlParams map[string]raml.NamedParameter) (params []definition.Parameter) {
	var sortedRamlParams []string
	for k := range ramlParams {
		sortedRamlParams = append(sortedRamlParams, k)
	}

	sort.Strings(sortedRamlParams)

	for _, name := range sortedRamlParams {
		ramlParam := ramlParams[name]

		param := new(definition.Parameter)

		// It takes the parameter name over the parameter key from raml definition
		if param.Name = name; ramlParam.Name != "" {
			param.Name = ramlParam.Name
		}

		param.Description = ramlParam.Description
		param.Type = tra.removeLibraryName(ramlParam.Type)
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
	if len(ramlHeaders) == 0 {
		return
	}

	var sortedRamlHeaders []string
	for k := range ramlHeaders {
		sortedRamlHeaders = append(sortedRamlHeaders, string(k))
	}

	sort.Strings(sortedRamlHeaders)

	for _, name := range sortedRamlHeaders {
		ramlHead := ramlHeaders[raml.HTTPHeader(name)]
		header := new(definition.Header)

		// It takes the parameter name over the parameter key from raml definition
		if header.Name = name; ramlHead.Name != "" {
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
	var sortedRamlTypes []string
	for k := range ramlTypes {
		sortedRamlTypes = append(sortedRamlTypes, string(k))
	}

	sort.Strings(sortedRamlTypes)

	for _, name := range sortedRamlTypes {
		ramlType := ramlTypes[name]
		customType := &definition.CustomType{
			Name:        name,
			Description: ramlType.Description,
			Type:        ramlType.Type,
			Enum:        ramlType.Enum,
			Default:     ramlType.Default,
		}

		// We need to remove the library name space from the type
		if t, ok := customType.Type.(string); ok {
			customType.Type = tra.removeLibraryName(t)
		}

		customType.Properties = tra.handleCustomTypeProperties(ramlType.Properties)

		for _, e := range ramlType.Examples {
			customType.Examples = append(customType.Examples, e)
		}

		if ramlType.Example != nil {
			customType.Examples = append(customType.Examples, ramlType.Example)
		}

		// It takes the parameter name over the parameter key from raml definition
		if ramlType.DisplayName != "" {
			customType.Name = ramlType.DisplayName
		}

		customTypes = append(customTypes, *customType)
	}

	return
}

// handleTraits Generic method which handles raml's trait definition.
func (tra *RamlTransformer) handleTraits(ramlTraits map[string]raml.Trait) (traits []definition.Trait) {
	var sortedRamlTraits []string
	for k := range ramlTraits {
		sortedRamlTraits = append(sortedRamlTraits, string(k))
	}

	sort.Strings(sortedRamlTraits)

	for _, name := range sortedRamlTraits {
		ramlTrait := ramlTraits[name]

		trait := definition.Trait{
			Name:        ramlTrait.Name,
			Usage:       ramlTrait.Usage,
			Description: ramlTrait.Description,
			Href: definition.Href{
				Parameters: tra.handleParameters(ramlTrait.QueryParameters),
			},
			Protocols: tra.handleProtocols(ramlTrait.Protocols),
		}

		trait.Transactions = tra.buildTransactions(ramlTrait.Headers, &ramlTrait.Bodies, ramlTrait.Responses)

		traits = append(traits, trait)
	}

	return
}

// handleSecuritySchemes Generic method which handles raml's security schemes definition.
func (tra *RamlTransformer) handleSecuritySchemes(ramlSchemes map[string]raml.SecurityScheme) (schemes []definition.SecurityScheme) {
	var sortedRamlSchemes []string
	for k := range ramlSchemes {
		sortedRamlSchemes = append(sortedRamlSchemes, string(k))
	}

	sort.Strings(sortedRamlSchemes)

	for _, name := range sortedRamlSchemes {
		ramlScheme := ramlSchemes[name]
		scheme := new(definition.SecurityScheme)

		// It takes the parameter name over the parameter key from raml definition
		if scheme.Name = name; ramlScheme.DisplayName != "" {
			scheme.Name = ramlScheme.DisplayName
		}

		scheme.Type = ramlScheme.Type
		scheme.Description = ramlScheme.Description

		var sortedRamlSchemesSets []string
		for k := range ramlScheme.Settings {
			sortedRamlSchemesSets = append(sortedRamlSchemesSets, string(k))
		}

		sort.Strings(sortedRamlSchemesSets)

		for _, k := range sortedRamlSchemesSets {
			scheme.Settings = append(scheme.Settings, definition.SecuritySchemeSetting{
				Name: k,
				Data: ramlScheme.Settings[k],
			})
		}

		scheme.Transactions = tra.buildTransactions(ramlScheme.DescribedBy.Headers, nil, ramlScheme.DescribedBy.Responses)

		schemes = append(schemes, *scheme)
	}

	return
}

// handleResources Generic method which handles raml's resources definition.
func (tra *RamlTransformer) handleResources(ramlResources interface{}, parent *definition.Resource) (resources []definition.Resource, err error) {
	switch rs := ramlResources.(type) {
	case map[string]raml.Resource:
		var sortedRs []string
		for k := range rs {
			sortedRs = append(sortedRs, k)
		}

		sort.Strings(sortedRs)

		for _, uri := range sortedRs {
			res := tra.handleResource(rs[uri], parent)
			if rs[uri].Nested != nil {
				if res.Resources, err = tra.handleResources(rs[uri].Nested, &res); err != nil {
					return
				}
			}
			resources = append(resources, res)
		}
	case map[string]*raml.Resource:
		var sortedRs []string
		for k := range rs {
			sortedRs = append(sortedRs, k)
		}

		sort.Strings(sortedRs)

		for _, uri := range sortedRs {
			res := tra.handleResource(*rs[uri], parent)
			if rs[uri].Nested != nil {
				if res.Resources, err = tra.handleResources(rs[uri].Nested, &res); err != nil {
					return
				}
			}
			resources = append(resources, res)
		}
	default:
		err = errors.New("The resource's type is unsupported")
	}
	return
}

// handleResource Generic method which handles raml's resource definition.
func (tra *RamlTransformer) handleResource(ramlRes raml.Resource, parent *definition.Resource) definition.Resource {
	return definition.Resource{
		Title:       ramlRes.DisplayName,
		Description: ramlRes.Description,
		Href: definition.Href{
			FullPath:   fmt.Sprintf("%s%s", parent.Href.FullPath, ramlRes.URI),
			Path:       ramlRes.URI,
			Parameters: tra.handleParameters(ramlRes.URIParameters),
		},
		Is:        tra.handleOptions(ramlRes.Is),
		SecuredBy: tra.handleOptions(ramlRes.SecuredBy),
		Actions:   tra.handleResourceMethods(ramlRes, ramlRes.Methods),
	}
}

// handleResource Generic method which handles raml's method definition.
func (tra *RamlTransformer) handleResourceMethods(ramlRes raml.Resource, ramlMethods []*raml.Method) (actions []definition.ResourceAction) {
	for _, ramlMethod := range ramlMethods {
		action := new(definition.ResourceAction)

		action.Title = ramlMethod.DisplayName
		action.Description = ramlMethod.Description
		action.Href = definition.Href{
			Parameters: tra.handleParameters(ramlMethod.QueryParameters),
		}
		action.Is = tra.handleOptions(ramlMethod.Is)

		// Inherits securedBy options from the parent securedBy if not present
		action.SecuredBy = tra.handleOptions(ramlMethod.SecuredBy)
		if action.SecuredBy == nil {
			action.SecuredBy = tra.handleOptions(ramlRes.SecuredBy)
		}

		action.Method = ramlMethod.Name
		action.Transactions = tra.buildTransactions(ramlMethod.Headers, &ramlMethod.Bodies, ramlMethod.Responses)

		actions = append(actions, *action)
	}

	return
}

// handleRequest It creates an API's request based on RAML's request parameters
func (tra *RamlTransformer) handleRequest(headers map[raml.HTTPHeader]raml.Header, bodies *raml.Bodies) (req *definition.Request) {
	if len(headers) == 0 && bodies == nil {
		return
	}

	return &definition.Request{
		Headers: tra.handleHeaders(headers),
		Body:    tra.handleBodies(bodies),
	}
}

// handleResponse It creates an API's response based on RAML's response
func (tra *RamlTransformer) handleResponse(code string, ramlResp raml.Response) (resp *definition.Response) {
	resp = new(definition.Response)

	resp.StatusCode, _ = strconv.Atoi(code)
	resp.Description = ramlResp.Description
	resp.Headers = tra.handleHeaders(ramlResp.Headers)
	resp.Body = tra.handleBodies(&ramlResp.Bodies)

	return
}

// handleBodies Generic method which handles raml's bodies definition.
func (tra *RamlTransformer) handleBodies(ramlBodies *raml.Bodies) (bodies []definition.Body) {
	if ramlBodies == nil || (ramlBodies.ApplicationJSON == nil && ramlBodies.Type == "") {
		return
	}

	body := new(definition.Body)

	if ramlBodies.ApplicationJSON != nil {
		body.MediaType = definition.MediaType("application/json")

		// t will be the body's type
		bodyType := tra.removeLibraryName(ramlBodies.ApplicationJSON.TypeString())

		// If properties is empty then it is not a api's CustomType
		if ramlBodies.ApplicationJSON.Properties != nil {
			customType := &definition.CustomType{
				Type: bodyType,
			}

			if ramlBodies.ApplicationJSON.Properties != nil {
				customType.Properties = tra.handleCustomTypeProperties(ramlBodies.ApplicationJSON.Properties)
			}

			body.CustomType = customType

		} else {
			body.Type = bodyType
		}

	} else {

		//@todo Add here your code to support other media types for bodies

		body.Type = tra.removeLibraryName(ramlBodies.Type)
		body.Description = ramlBodies.Description
		body.Example = ramlBodies.Example

		if ramlBodies.Example != "" {
			body.Example = ramlBodies.Example
		}
	}

	bodies = append(bodies, *body)
	return
}

// handleCustomTypeProperties It transforms RAML's custom properties into an API's array of definition.property
func (tra *RamlTransformer) handleCustomTypeProperties(properties map[string]interface{}) (props []definition.CustomTypeProperty) {
	var sortedProps []string
	for k := range properties {
		sortedProps = append(sortedProps, k)
	}

	sort.Strings(sortedProps)

	for _, k := range sortedProps {
		p := properties[k]
		// convert from map of interface to property
		mapToProperty := func(val map[interface{}]interface{}) definition.CustomTypeProperty {
			var p definition.CustomTypeProperty
			p.Required = true
			for k, v := range val {
				switch k {
				case "type":
					p.Type = tra.removeLibraryName(v.(string))
				case "required":
					p.Required = v.(bool)
				case "description":
					p.Description = v.(string)
				case "example":
					p.Example = v.(string)
				case "properties":
					if properties, ok := v.(map[interface{}]interface{}); ok {
						props := make(map[string]interface{})
						for name, prop := range properties {
							props[name.(string)] = prop
						}
						p.Properties = tra.handleCustomTypeProperties(props)
					}
				}
			}
			return p
		}

		prop := definition.CustomTypeProperty{Required: true}
		switch p.(type) {
		case string:
			prop.Type = tra.removeLibraryName(p.(string))
		case map[interface{}]interface{}:
			prop = mapToProperty(p.(map[interface{}]interface{}))
		case definition.CustomTypeProperty:
			prop = p.(definition.CustomTypeProperty)
		}

		if prop.Type == "" { // if has no type, we set it as string
			prop.Type = "string"
		}

		prop.Name = k

		// if has "?" suffix, remove the "?" and set required=false
		if strings.HasSuffix(prop.Name, "?") {
			prop.Required = false
			prop.Name = prop.Name[:len(prop.Name)-1]
		}

		props = append(props, prop)
	}

	return
}

// buildTransactions It holds the responsibility to create multiple transactions based on RAML's request/responses
func (tra *RamlTransformer) buildTransactions(headers map[raml.HTTPHeader]raml.Header, bodies *raml.Bodies, responses map[raml.HTTPCode]raml.Response) (transactions []definition.Transaction) {
	req := tra.handleRequest(headers, bodies)

	if len(responses) > 0 {
		var sortedResponses []string
		for k := range responses {
			sortedResponses = append(sortedResponses, string(k))
		}

		sort.Strings(sortedResponses)

		for _, code := range sortedResponses {
			ramlResp := responses[raml.HTTPCode(code)]

			resp := tra.handleResponse(code, ramlResp)

			trans := definition.Transaction{
				Response: *resp,
			}

			if req != nil {
				trans.Request = *req
			}

			// Discard the request for the next iterations since it will be duplicated for each transaction-request
			req = nil

			transactions = append(transactions, trans)
		}
	} else if req != nil {
		transactions = append(transactions, definition.Transaction{Request: *req})
	}

	return
}

// removeLibraryName It removes the library's namespace from a string
func (tra *RamlTransformer) removeLibraryName(name string) string {
	s := strings.Split(strings.TrimSpace(name), ".")
	if len(s) == 2 {
		return s[1]
	}
	return name
}
