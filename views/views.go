package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

//Init loads templates by pairing the layout templates with the normal ones
func Init(templatesDir string) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	//read all html files from the templatesDir
	templateFiles, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layout and templates/ directory
	for _, templateFile := range templateFiles {
		if templateFile != "templates/layout.html" {
			templates[filepath.Base(templateFile)] = template.Must(template.ParseFiles(templateFile, "templates/layout.html"))
		}
	}
}

//RenderTemplate renders a template by name into a ResponseWriter
//returns error if the template cannot be found by that name
func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	fmt.Println("templateName", templateName)
	tmpl, ok := templates[templateName]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", templateName)
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return err
	}

	return nil
}
