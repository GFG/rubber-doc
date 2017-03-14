package transformer

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/walker"
)

type BlueprintTransformer struct{}

func NewBlueprintTransformer() Transformer {
	return new(BlueprintTransformer)
}

func (f *BlueprintTransformer) Transform(data interface{}) (def *definition.Api, err error) {
	el, ok := data.(walker.ObjectWalker)
	if !ok {
		err = errors.New("The data's struct given isn't supported by the Blueprint's Transformer")
		return
	}

	def = f.handleContent(&el)

	return
}

//Blueprint transformation methods, names are bound to what they parse

func (f *BlueprintTransformer) resourceGroups(el *walker.ObjectWalker, apiDef *definition.Api) {
	children := filterContentByClass("resourceGroup", el)

	for _, child := range children {
		g := &definition.ResourceGroup{
			Title:       child.Path("meta.title").String(),
			Description: f.handleDescription(child),
		}

		f.resources(child, g)
		apiDef.ResourceGroups = append(apiDef.ResourceGroups, *g)
	}
}

func (f *BlueprintTransformer) resources(el *walker.ObjectWalker, g *definition.ResourceGroup) {
	children := filterContentByElement("resource", el)

	cr := make(chan *definition.Resource)
	oc := make([]string, len(children))
	rs := make([]definition.Resource, len(children))

	for i, child := range children {
		oc[i] = child.Path("meta.title").String()

		go func(c *walker.ObjectWalker) {
			r := &definition.Resource{
				Title:       c.Path("meta.title").String(),
				Description: f.handleDescription(c),
				Href:        f.handleHref(c),
			}

			f.resourceAction(c, r)

			cr <- r
		}(child)
	}

	for i := 0; i < len(children); i++ {
		r := <-cr

		for n := range oc {
			if oc[n] == r.Title {
				rs[n] = *r
			}
		}
	}

	g.Resources = rs
}

func (f *BlueprintTransformer) resourceAction(el *walker.ObjectWalker, r *definition.Resource) {
	children := filterContentByElement("transition", el)

	for _, child := range children {
		t := &definition.ResourceAction{
			Title:       child.Path("meta.title").String(),
			Description: f.handleDescription(child),
			Href:        f.handleHref(child),
		}

		f.transactions(child, t)
		r.Actions = append(r.Actions, *t)
	}
}

func (f *BlueprintTransformer) transactions(el *walker.ObjectWalker, ra *definition.ResourceAction) {
	children := filterContentByElement("httpTransaction", el)

	for _, child := range children {
		cx, err := child.Path("content").Children()
		if err != nil {
			continue
		}

		ra.Transactions = append(ra.Transactions, f.transaction(cx))
	}
}

func (f *BlueprintTransformer) transaction(el []*walker.ObjectWalker) definition.Transaction {
	x := new(definition.Transaction)

	for _, child := range el {
		if child.Path("element").String() == "httpRequest" {
			f.request(child, x)
		}

		if child.Path("element").String() == "httpResponse" {
			f.response(child, x)
		}
	}

	return *x
}

func (f *BlueprintTransformer) request(child *walker.ObjectWalker, x *definition.Transaction) {
	x.Request.Title = child.Path("meta.title").String()
	x.Request.Description = f.handleDescription(child)
	x.Request.Headers = f.handleHeaders(child.Path("attributes.headers"))

	cx, err := child.Path("content").Children()
	if err != nil {
		return
	}

	for _, c := range cx {
		if hasClass("messageBody", c) {
			x.Request.Body = []definition.Body{
				f.handleBody(c),
			}
		}
	}
}

func (f *BlueprintTransformer) response(child *walker.ObjectWalker, x *definition.Transaction) {
	s := child.Path("attributes.statusCode").String()
	n, err := strconv.Atoi(s)
	if err != nil {
		n = 0
	}

	x.Response.StatusCode = n
	x.Response.Headers = f.handleHeaders(child.Path("attributes.headers"))
	x.Response.Description = f.handleDescription(child)

	cx, err := child.Path("content").Children()
	if err != nil {
		return
	}

	for _, c := range cx {
		if hasClass("messageBody", c) {
			x.Response.Body = []definition.Body{
				f.handleBody(c),
			}
		}

	}
}

//Internals methods, handle specific nodes & content retrieval

//Handle the content to start parsing the document elemens
func (f *BlueprintTransformer) handleContent(el *walker.ObjectWalker) *definition.Api {
	apiDef := &definition.Api{}

	children, _ := el.Path("content").Children()
	for _, child := range children {
		f.handleElements(child, apiDef)
	}

	return apiDef
}

//Handle metadata extraction
func (f *BlueprintTransformer) handleMetadata(el *walker.ObjectWalker) (m map[string]string) {
	children, err := el.Path("attributes.meta").Children()
	m = make(map[string]string)

	if err != nil {
		return
	}

	for _, v := range children {
		m[strings.ToLower(v.Path("content.key.content").String())] = v.Path("content.value.content").String()
	}

	return
}

//Handle element parsing
func (f *BlueprintTransformer) handleElements(el *walker.ObjectWalker, apiDef *definition.Api) {
	switch el.Path("element").String() {
	case "category":
		if hasClass("api", el) {
			f.handleTitles(el, apiDef)

			meta := f.handleMetadata(el)
			if version, ok := meta["version"]; ok {
				apiDef.Version = version
			}

			if host, ok := meta["host"]; ok {
				apiDef.BaseURI = host
				proto, _ := definition.NewProtocolFromURL(host)
				apiDef.Protocols = append(apiDef.Protocols, proto)
			}

			f.resourceGroups(el, apiDef)
		}
	}
}

//Handle the titles
func (f *BlueprintTransformer) handleTitles(el *walker.ObjectWalker, apiDef *definition.Api) {
	if hasClass("api", el) {
		apiDef.Title = el.Path("meta.title").String()
	}
}

//Handle the descriptions section (named copy with APIB)
func (f *BlueprintTransformer) handleDescription(el *walker.ObjectWalker) string {
	children, err := el.Path("content").Children()
	if err != nil {
		return ""
	}

	for _, child := range children {
		if child.Path("element").String() == "copy" {
			return child.Path("content").String()
		}
	}

	return ""
}

//Handle the href sections, including it's internal params
func (f *BlueprintTransformer) handleHref(child *walker.ObjectWalker) (h definition.Href) {
	href := child.Path("attributes.href")

	if href.Value().IsValid() {
		h.Path = href.String()
	}

	contents, err := child.Path("attributes.hrefVariables.content").Children()
	if err != nil {
		return
	}

	for _, content := range contents {
		v := &definition.Parameter{
			Required:    contains("attributes.typeAttributes", "required", content),
			Type:        content.Path("content.value.element").String(),
			Example:     content.Path("content.value.content").String(),
			Name:        content.Path("content.key.content").String(),
			Description: content.Path("meta.description").String(),
		}

		h.Parameters = append(h.Parameters, *v)
	}

	return
}

//Handle body examples, both for response and request
func (f *BlueprintTransformer) handleBody(child *walker.ObjectWalker) (a definition.Body) {
	if child.Path("element").String() == "asset" {
		ms := map[string]string{
			`\\n`: `\n`,
			`\\r`: `\r`,
			`\\"`: `\"`,
		}
		content := child.Path("content").String()

		for key, val := range ms {
			content = strings.Replace(content, key, val, -1)
		}

		return definition.Body{
			MediaType: definition.MediaType(child.Path("attributes.contentType").String()),
			Example:   content,
		}
	}

	return
}

//Handle the header extraction
func (f *BlueprintTransformer) handleHeaders(child *walker.ObjectWalker) (hs []definition.Header) {
	if child.Path("element").String() == "httpHeaders" {
		contents, err := child.Path("content").Children()
		if err != nil {
			return
		}

		for _, content := range contents {
			h := definition.Header{
				Name:    content.Path("content.key.content").String(),
				Example: content.Path("content.value.content").String(),
			}

			hs = append(hs, h)
		}

		return
	}

	return
}

//generic helper methods, unbound to struct

func hasClass(s string, child *walker.ObjectWalker) bool {
	return contains("meta.classes", s, child)
}

func contains(key, s string, child *walker.ObjectWalker) bool {
	v := child.Path(key).Value()

	if !v.IsValid() {
		return false
	}

	for i := 0; i < v.Len(); i++ {
		if s == v.Index(i).Interface().(string) {
			return true
		}
	}

	return false
}

func filterContentByClass(s string, el *walker.ObjectWalker) (xs []*walker.ObjectWalker) {
	children, err := el.Path("content").Children()
	if err != nil {
		return
	}

	for _, child := range children {
		if hasClass(s, child) {
			xs = append(xs, child)
		}
	}

	return
}

func filterContentByElement(s string, el *walker.ObjectWalker) (xs []*walker.ObjectWalker) {
	children, err := el.Path("content").Children()
	if err != nil {
		return
	}

	for _, child := range children {
		if child.Path("element").String() == s {
			xs = append(xs, child)
		}
	}

	return
}
