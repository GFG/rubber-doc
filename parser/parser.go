package parser

import (
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
)

// Parser
type Parser interface {
	Parse(filename string, tra transformer.Transformer) (def *definition.Api, err error)
}
