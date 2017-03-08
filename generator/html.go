package generator

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator/config"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/generator/html"
)

// HTML Represents a html's generator
type HTML struct {
	templates []*html.Template
}

// NewHTMLGenerator Returns a HTMLGenerator's struct
func NewHTMLGenerator(cfg config.Config, data definition.Api) (gen Generator, err error) {
	htmlGen := new(HTML)

	if cfg.IsCombined() {
		err = htmlGen.populateWithCombinedTemplates(cfg, data)
	} else {
		err = htmlGen.populateWithTemplates(cfg, data)
	}

	if err == nil {
		gen = htmlGen
	}

	return
}

// Generate Generates the output based on the templates given
func (gen *HTML) Generate() (err error) {
	if gen.templates == nil {
		err = errors.New("There is no templates to be processed by the HTML's generator.")
	}

	for _, tmpl := range gen.templates {
		err = tmpl.Execute()
	}
	return
}

// populateWithCombinedTemplates It's responsible to collect all templates from the configuration and create one HTML's template
func (gen *HTML) populateWithCombinedTemplates(cfg config.Config, data definition.Api) (err error) {
	var (
		filenames []string
		template  *html.Template
	)

	for _, tc := range cfg.Templates() {
		filenames = append(filenames, filepath.Join(cfg.Src(), tc.Src()))
	}

	// For combined, the first template's filename will be used as template's name
	name := filepath.Base(filenames[0])
	output := cfg.Output()

	if template, err = html.NewTemplate(name, data, filenames, output); err != nil {
		return
	}

	gen.templates = append(gen.templates, template)

	return
}

// populateWithTemplates It's responsible to create one HTML's template by each template defined on the configuration
func (gen *HTML) populateWithTemplates(cfg config.Config, data definition.Api) (err error) {
	var (
		filenames []string
		template  *html.Template
	)

	for _, tc := range cfg.Templates() {
		filenames = append(filenames, filepath.Join(cfg.Src(), tc.Src()))
		name := filepath.Base(filenames[0])
		output := filepath.Join(cfg.Dst(), tc.Dst())

		if template, err = html.NewTemplate(name, data, filenames, output); err != nil {
			return
		}

		gen.templates = append(gen.templates, template)
	}

	return
}
