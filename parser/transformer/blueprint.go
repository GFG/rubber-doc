package transformer

import (
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/walker"
)

type BlueprintTransformer struct {
	def *definition.Api
}

func NewBlueprintTransformer() Transformer {
	return new(BlueprintTransformer)
}

func (f BlueprintTransformer) Transform(data interface{}) (def *definition.Api) {
	el, ok := data.(walker.ObjectWalker)
	if !ok {
		return
	}

	f.def = &definition.Api{}

	f.digContent(&el)

	return f.def
}

func (f *BlueprintTransformer) digContent(el *walker.ObjectWalker) {
	children, _ := el.Path("content").Children()
	for _, child := range children {
		f.digElements(child)
	}
}

func (f *BlueprintTransformer) digTitle(el *walker.ObjectWalker) {
	if hasClass("api", el) {
		f.def.Title = el.Path("meta.title").String()
	}
}

func (f *BlueprintTransformer) digElements(el *walker.ObjectWalker) {
	switch el.Path("element").String() {
	case "category":
		if hasClass("api", el) {
			f.digTitle(el)
		}
	}
}

func hasClass(s string, child *walker.ObjectWalker) bool {
	return isContains("meta.classes", s, child)
}

func isContains(key, s string, child *walker.ObjectWalker) bool {
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
