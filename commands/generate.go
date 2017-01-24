package commands

import (
	"fmt"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser"
)

const (
	BLUEPRINT = "blueprint"
	RAML = "raml"
)

type GenerateCommand struct {
	InputFormat  string
	OutputFormat string
	Src          string
	OutputDir    string
}

func (c *GenerateCommand) Execute() (err error) {
	if c.InputFormat == BLUEPRINT {
		var spec parser.Specification
		bpParser := parser.NewBlueprintParser()

		if spec, err = bpParser.Parse(c.Src, new(parser.BlueprintFormatter)); err != nil {
			return
		}

		fmt.Printf("%+v\n", spec)
	}

	if c.InputFormat == RAML {
		var spec parser.Specification
		rParser := parser.NewRamlParser()

		if spec, err = rParser.Parse(c.Src, new(parser.RamlFormatter)); err != nil {
			return
		}

		fmt.Printf("%+v\n", spec)
	}

	return
}

