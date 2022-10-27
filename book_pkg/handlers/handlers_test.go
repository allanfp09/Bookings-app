package handlers

import (
	"bookings/book_pkg/models"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type postData struct {
	key   string
	value string
}

var testData = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"reservation", "/reservation", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	route := getRoutes()
	ts := httptest.NewTLSServer(route)
	defer ts.Close()

	for _, td := range testData {
		if td.method == "GET" {
			req, err := ts.Client().Get(ts.URL + td.url)
			if err != nil {
				t.Fatal(err)
			}

			if req.StatusCode != td.expectedStatusCode {
				t.Errorf("for %s, expected status code %d, but got %d instead", td.name, td.expectedStatusCode, req.StatusCode)
			}
		} else {
			values := url.Values{}

			resp, err := ts.Client().PostForm(ts.URL+td.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != td.expectedStatusCode {
				t.Errorf("for %s, expected %d, got %d", td.name, td.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservations{
		RoomID: 1,
		Room: models.Rooms{
			ID:       1,
			RoomName: "General's Room",
		},
	}

	r, err := http.NewRequest("GET", "/reservation", nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := getCtx(r)
	r = r.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// When reservation session has no data
	r, err = http.NewRequest("GET", "/reservation", nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test when room id is incorrect
	r, err = http.NewRequest("GET", "/reservation", nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 8
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2022-10-02")
	ed, _ := time.Parse(layout, "2022-10-05")

	reservation := models.Reservations{
		StartDate: sd,
		EndDate:   ed,
		RoomID:    1,
		Room: models.Rooms{
			ID:       1,
			RoomName: "General's Room",
		},
	}

	//rBody := "start_date=2050-01-01"
	//rBody = fmt.Sprintf("%s&%s", rBody, "end_date=2050-01-01")
	//rBody = fmt.Sprintf("%s&%s", rBody, "first_name=Allan")
	//rBody = fmt.Sprintf("%s&%s", rBody, "last_name=Fuentes")
	//rBody = fmt.Sprintf("%s&%s", rBody, "email=allan@allan.com")
	//rBody = fmt.Sprintf("%s&%s", rBody, "phone=88193344")
	//rBody = fmt.Sprintf("%s&%s", rBody, "room_id=1")

	form := url.Values{}
	form.Add("first_name", "allan")
	form.Add("last_name", "fuentes")
	form.Add("email", "allan@allan.com")
	form.Add("phone", "8383923")
	form.Add("room_id", "1")

	r, _ := http.NewRequest("POST", "/reservation", strings.NewReader(form.Encode()))
	ctx := getCtx(r)
	r = r.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing post body
	r, _ = http.NewRequest("POST", "/reservation", nil)
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test form is invalid
	form = url.Values{}
	form.Add("first_name", "a")
	form.Add("last_name", "f")
	form.Add("email", "allan@allan")
	form.Add("phone", "8383923")
	form.Add("room_id", "1")

	r, _ = http.NewRequest("POST", "/reservation", strings.NewReader(form.Encode()))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// Test when session is not set with reservation

	r, _ = http.NewRequest("POST", "/reservation", nil)
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)

	}

	// Test when unable to insert reservation
	form = url.Values{}
	form.Add("first_name", "allan")
	form.Add("last_name", "fuentes")
	form.Add("email", "allan@allan")
	form.Add("phone", "8383923")
	form.Add("start_date", "2022-10-25")
	form.Add("end_date", "2022-10-27")
	form.Add("room_id", "1")

	r, _ = http.NewRequest("POST", "/reservation", strings.NewReader(form.Encode()))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	reservation.RoomID = 2

	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test when unable to insert room restriction
	form = url.Values{}
	form.Add("first_name", "allan")
	form.Add("last_name", "fuentes")
	form.Add("email", "allan@allan")
	form.Add("phone", "8383923")
	form.Add("start_date", "2022-10-25")
	form.Add("end_date", "2022-10-27")
	form.Add("room_id", "1")

	r, _ = http.NewRequest("POST", "/reservation", strings.NewReader(form.Encode()))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	reservation.RoomID = 1000

	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
