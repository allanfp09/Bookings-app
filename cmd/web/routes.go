package main

import (
	"bookings/book_pkg/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func routes() http.Handler {
	r := mux.NewRouter()

	//r.Use(CsrfToken)
	r.Use(LoadSession)

	r.HandleFunc("/", handlers.Repo.Home).Methods("GET")
	r.HandleFunc("/about", handlers.Repo.About).Methods("GET")
	r.HandleFunc("/search-availability", handlers.Repo.Availability).Methods("GET")
	r.HandleFunc("/search-availability", handlers.Repo.PostAvailability).Methods("POST")
	r.HandleFunc("/choose-room/{id}", handlers.Repo.ChooseRoom).Methods("GET")
	r.HandleFunc("/reservation", handlers.Repo.Reservation).Methods("GET")
	r.HandleFunc("/reservation", handlers.Repo.PostReservation).Methods("POST")
	r.HandleFunc("/reservation-summary", handlers.Repo.ReservationSummary).Methods("GET")

	r.PathPrefix("static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return r
}
