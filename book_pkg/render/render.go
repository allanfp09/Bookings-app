package render

import (
	"bookings/book_pkg/config"
	"bookings/book_pkg/models"
	"fmt"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var tmplPath = "./templates"
var functions = template.FuncMap{}
var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddUsefulData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.Csrf = nosurf.Token(r)

	return data
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc = make(map[string]*template.Template)
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = TemplateCache()
	}

	t := tc[tmpl]

	data := AddUsefulData(td, r)

	err := t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func TemplateCache() (map[string]*template.Template, error) {
	var tmplCache = make(map[string]*template.Template)

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", tmplPath))
	if err != nil {
		return tmplCache, nil
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return tmplCache, nil
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))
		if err != nil {
			return tmplCache, nil
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))
			if err != nil {
				return tmplCache, nil
			}
		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}
