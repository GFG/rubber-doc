package command

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
)

const (
	BLUEPRINT = ".apib"
	RAML      = ".raml"
)

type GenerateCommand struct {
	OutputFormat string
	Src          string
	OutputDir    string
}

func (c *GenerateCommand) Execute() error {
	var (
		p parser.Parser
		f transformer.Transformer
	)

	switch filepath.Ext(c.Src) {
	case BLUEPRINT:
		p = parser.NewBlueprintParser()
		f = transformer.NewBlueprintTransformer()
	case RAML:
		p = parser.NewRamlParser()
		f = transformer.NewRamlTransformer()
	default:
		return errors.New("unsupported format")
	}

	def, err := p.Parse(c.Src, f)
	if def != nil {
		fmt.Printf("%+v\n", def)
	}

	return err
}
