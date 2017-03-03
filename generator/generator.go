package generator

import (
	"html/template"
	"os"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
)

type Generator struct {
	def definition.Api
	cfg Config
}

// NewGenerator Returns a Generator's struct
func NewGenerator(configFile string, def definition.Api) (gen *Generator, err error) {
	var cfg Config
	if cfg, err = NewConfigYaml(configFile); err == nil {
		gen = &Generator{def, cfg}
	}
	return
}

// Generate Generates the output based on the templates given
func (gen *Generator) Generate() (err error) {
	var t *template.Template

	if gen.config().IsCombined() {
		var tmpls []string
		for _, tmpl := range gen.config().TemplatesConfig() {
			tmpls = append(tmpls, tmpl.Src())
		}

		if t, err = template.ParseFiles(tmpls...); err == nil {
			err = gen.processTemplate(t, gen.config().Output())
		}

	} else {
		for _, tmpl := range gen.config().TemplatesConfig() {
			if t, err = template.ParseFiles(tmpl.Src()); err == nil {
				err = gen.processTemplate(t, tmpl.Dst())
			}
		}
	}
	return
}

//processTemplate Parses the templates and creates the output.
func (gen *Generator) processTemplate(t *template.Template, outputFile string) (err error) {
	if err = createDir(gen.config().Dst()); err != nil {
		return
	}

	var f *os.File
	if f, err = os.Create(outputFile); err != nil {
		return
	}
	defer f.Close()

	err = t.Execute(f, gen.data())

	return
}

// Config Returns the generator's configuration
func (gen *Generator) config() Config {
	return gen.cfg
}

// data Returns the generator's data
func (gen *Generator) data() interface{} {
	return gen.def
}

// createDir Creates a directory if not exist
func createDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return
}
