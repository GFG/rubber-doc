package html

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/definition"
)

// Template Represents html's template handler
type Template struct {
	handler *template.Template
	data    definition.Api
	name    string
	output  string
}

// NewTemplate
func NewTemplate(name string, data definition.Api, filenames []string, output string) (tmpl *Template, err error) {
	handler := template.New(name)

	handler.Funcs(helpers(data))

	if handler, err = handler.ParseFiles(filenames...); err == nil {
		tmpl = &Template{handler, data, name, output}
	}
	return
}

// Execute Parses the templates and creates the output.
func (t *Template) Execute() (err error) {
	if err = createDir(t.output); err != nil {
		return
	}

	var f *os.File
	if f, err = os.Create(t.output); err != nil {
		return
	}
	defer f.Close()

	err = t.handler.Execute(f, t.data)

	return
}

// helpers Returns the helpers given to the templates
func helpers(data definition.Api) template.FuncMap {
	return template.FuncMap{
		"Comment": func(t string) template.HTML {
			return template.HTML(t)
		},
		"Lower":            strings.ToLower,
		"CustomTypeByName": data.CustomTypeByName,
	}
}

// createDir Creates a directory if not exist
func createDir(filename string) (err error) {
	dir := filepath.Dir(filename)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	return
}
