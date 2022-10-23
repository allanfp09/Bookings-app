package models

import "bookings/book_pkg/forms"

type TemplateData struct {
	Csrf      string
	Forms     *forms.Forms
	Data      map[string]interface{}
	StringMap map[string]string
}
