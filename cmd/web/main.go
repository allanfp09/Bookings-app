package main

import (
	"bookings/book_pkg/config"
	"bookings/book_pkg/driver"
	"bookings/book_pkg/handlers"
	"bookings/book_pkg/helpers"
	"bookings/book_pkg/models"
	"bookings/book_pkg/render"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Application listening to localhost%s", PORT))
	server := &http.Server{
		Addr:    PORT,
		Handler: routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	gob.Register(models.Reservations{})
	gob.Register(models.Rooms{})
	gob.Register(models.RoomRestrictions{})
	gob.Register(models.Restrictions{})
	gob.Register(models.Users{})

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

	// Connecting to a postgres database
	log.Println("Connecting to the database...")
	db, err := driver.ConnectDBS("postgres://allanfp:fuenpi09@localhost:5432/bookings")
	if err != nil {
		log.Fatal("cannot connect to the database")
	}

	// Getting template cache
	tc, err := render.TemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)

	render.NewRenderer(&app)

	return db, nil
}
