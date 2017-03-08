package command

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator/config"
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
		p     parser.Parser
		trans transformer.Transformer
	)

	format := filepath.Ext(c.SpecFile)
	switch format {
	case BLUEPRINT:
		p = parser.NewBlueprintParser()
		trans = transformer.NewBlueprintTransformer()
	case RAML:
		p = parser.NewRamlParser()
		trans = transformer.NewRamlTransformer()
	default:
		err = errors.Wrapf(err, "The format found %s for the specification given is unsuported", format)
		return
	}

	var def *definition.Api
	if def, err = p.Parse(c.SpecFile, trans); err != nil {
		return
	}

	var cfg config.Config
	if cfg, err = config.FromYaml(c.ConfigFile); err != nil {
		return
	}

	var gen generator.Generator
	gen, err = generator.NewHTMLGenerator(cfg, *def)

	if err != nil {
		return
	}

	err = gen.Generate()

	return
}
