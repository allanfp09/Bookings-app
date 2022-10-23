package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

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
