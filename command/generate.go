package command

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser/transformer"
)

const (
	BLUEPRINT = ".apib"
	RAML      = ".raml"
)

// GenerateCommand Represents the struct of the generate command
type GenerateCommand struct {
	SpecFile   string
	ConfigFile string
}

// Execute
func (c *GenerateCommand) Execute() (err error) {
	var (
		p parser.Parser
		f transformer.Transformer
	)

	format := filepath.Ext(c.SpecFile)
	switch format {
	case BLUEPRINT:
		p = parser.NewBlueprintParser()
		f = transformer.NewBlueprintTransformer()
	case RAML:
		p = parser.NewRamlParser()
		f = transformer.NewRamlTransformer()
	default:
		err = errors.Wrapf(err, "The format found %s for the specification given is unsuported", format)
		return
	}

	var def *definition.Api
	if def, err = p.Parse(c.SpecFile, f); err != nil {
		return
	}

	var gen *generator.Generator
	if gen, err = generator.NewGenerator(c.ConfigFile, *def); err != nil {
		return
	}

	err = gen.Generate()

	return
}
