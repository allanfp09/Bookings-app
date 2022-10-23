package render

import (
	"bookings/book_pkg/models"
	"net/http"
	"testing"
)

func TestTemplate(t *testing.T) {
	tmplPath = "./../../templates"
	r, err := http.NewRequest("GET", "/whatever", nil)
	if err != nil {
		t.Fatal(err)
	}
	var mw myWriter
	tc, err := TemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	app.TemplateCache = tc

	err = Template(&mw, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing to template")
	}
}

func TestTemplateCache(t *testing.T) {
	_, err := TemplateCache()
	if err != nil {
		t.Error("no template cache found")
	}
}

func TestNewTemplate(t *testing.T) {
	NewTemplate(app)
}
