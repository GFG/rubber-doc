package transformer

import "github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"

type Transformer interface {
	Transform(data interface{}) *definition.Api
}
