package handlers

import (
	"bookings/book_pkg/config"
	"bookings/book_pkg/helpers"
	"bookings/book_pkg/models"
	"bookings/book_pkg/render"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var tmplPath = "./../../templates"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	gob.Register(models.Reservations{})

	// Formatting errors
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	app.ErrorLog = errorLog

	infoLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	helpers.NewHelpers(&app)

	// App configurations
	app.InProduction = false
	app.UseCache = true

	// Initializing a session
	session = scs.New()
	session.Cookie.Path = "/"
	session.Cookie.Secure = app.InProduction
	session.Cookie.HttpOnly = true
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	app.Sessions = session

	// Getting template cache
	tc, err := TemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc

	repo := NewTestRepo(&app)
	NewHandler(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	gob.Register(models.Reservations{})

	// Formatting errors
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	app.ErrorLog = errorLog

	infoLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	helpers.NewHelpers(&app)

	// App configurations
	app.InProduction = false
	app.UseCache = true

	// Initializing a session
	session = scs.New()
	session.Cookie.Path = "/"
	session.Cookie.Secure = app.InProduction
	session.Cookie.HttpOnly = true
	session.Lifetime = time.Hour * 24
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	app.Sessions = session

	// Getting template cache
	tc, err := TemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc

	repo := NewTestRepo(&app)
	NewHandler(repo)

	render.NewRenderer(&app)

	r := mux.NewRouter()

	//r.Use(CsrfToken)
	r.Use(LoadSession)

	r.HandleFunc("/", Repo.Home).Methods("GET")
	r.HandleFunc("/about", Repo.About).Methods("GET")
	r.HandleFunc("/reservation", Repo.Reservation).Methods("GET")
	r.HandleFunc("/reservation", Repo.PostReservation).Methods("POST")
	r.HandleFunc("/reservation-summary", Repo.ReservationSummary).Methods("GET")

	r.PathPrefix("static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return r
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

func CsrfToken(next http.Handler) http.Handler {
	csrf := nosurf.New(next)
	csrf.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrf
}

func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
