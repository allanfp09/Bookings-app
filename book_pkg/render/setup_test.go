package render

import (
	"bookings/book_pkg/config"
	"bookings/models"
	"encoding/gob"
	"net/http"
	"os"
	"testing"
)

var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	testApp.InProduction = false
	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWriter) Write(b []byte) (int, error) {
	v := len(b)
	return v, nil
}

func (mw *myWriter) WriteHeader(statusCode int) {}
