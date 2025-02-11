package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
)

var TemplateMap map[string]*template.Template

func LoadTemplates(pattern string) error {
	tpls, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}

	TemplateMap = make(map[string]*template.Template)
	for _, tpl := range tpls.Templates() {
		TemplateMap[tpl.Name()] = tpl
		// Log only if the template name doesn't end with ".html"
		if !strings.HasSuffix(tpl.Name(), ".html") {
			log.Printf("Loaded template: %s", tpl.Name())
		}
	}
	log.Println("tmpl map: ", TemplateMap)
	return nil
}

func RenderTemplate(w io.Writer, name string, data interface{}) error {
	tmpl, ok := TemplateMap[name]
	if !ok {
		return fmt.Errorf("template %q not found", name)
	}
	return tmpl.Execute(w, data)
}
