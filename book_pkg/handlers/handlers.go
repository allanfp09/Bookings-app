package handlers

import (
	"bookings/book_pkg/config"
	"bookings/book_pkg/driver"
	"bookings/book_pkg/forms"
	"bookings/book_pkg/helpers"
	"bookings/book_pkg/models"
	"bookings/book_pkg/render"
	"bookings/book_pkg/repository"
	"bookings/book_pkg/repository/dbrepo"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

func NewTestRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewTestingRepo(app),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := repo.App.Sessions.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	repo.App.Sessions.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	var data = make(map[string]interface{})
	data["reservation"] = res
	_ = render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
		Forms:     forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Sessions.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.RequiredField("first_name", "last_name", "email")
	form.FieldLength("first_name", 3)
	form.IsEmail("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		_ = render.Template(w, r, "reservation.page.tmpl", &models.TemplateData{
			Forms: form,
			Data:  data,
		})
		return
	}

	reservationID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestrictions{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: reservationID,
		RestrictionID: 1,
	}
	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	repo.App.Sessions.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	var data = make(map[string]interface{})
	reservation, ok := repo.App.Sessions.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		repo.App.InfoLog.Println("could not get any data from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	repo.App.Sessions.Remove(r.Context(), "reservation")
	data["reservation"] = reservation

	_ = render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		repo.App.InfoLog.Println("no rooms")
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		repo.App.InfoLog.Println("no rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	var data = make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservations{
		StartDate: startDate,
		EndDate:   endDate,
	}

	repo.App.Sessions.Put(r.Context(), "reservation", res)

	_ = render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	roomID, err := strconv.Atoi(param["id"])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := repo.App.Sessions.Get(r.Context(), "reservation").(models.Reservations)
	if !ok {
		helpers.ServerError(w, errors.New("no reservation session found"))
		return
	}

	res.RoomID = roomID

	repo.App.Sessions.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservation", http.StatusSeeOther)
}
