package commands

import (
	"fmt"
	"os"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/parser"
)

const BLUEPRINT = "blueprint"

type GenerateCommand struct {
	InputFormat  string
	OutputFormat string
	Src          string
	OutputDir    string
}

func (c *GenerateCommand) Execute() error {
	if c.InputFormat == BLUEPRINT {
		f, err := os.Open(c.Src)

		if err != nil {
			return fmt.Errorf("Failed to open provided file: %s", c.Src)
		}

		defer f.Close()
		bpParser := parser.Blueprint{}
		bp, err := bpParser.Parse(f)

		if err != nil {
			return err
		}

		bpParser.PrintRecursiveMap(bp)
	}

	return nil
}
