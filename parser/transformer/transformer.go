package transformer

import "github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"

type Transformer interface {
	Transform(data interface{}) (def *definition.Api, err error)
}
