package forms

import (
	govalid "github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

type Forms struct {
	url.Values
	Errors errors
}

func New(val url.Values) *Forms {
	return &Forms{
		val,
		errors(map[string][]string{}),
	}
}

func (f *Forms) IsValid() bool {
	return len(f.Errors) == 0
}

func (f *Forms) RequiredField(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Forms) FieldLength(field string, length int) bool {
	fd := f.Get(field)
	if len(fd) < length {
		f.Errors.Add(field, "length must be longer")
		return true
	}

	return false
}

func (f *Forms) IsEmail(field string) {
	value := f.Get(field)
	if !govalid.IsEmail(value) {
		f.Errors.Add(field, "email is not correct")
	}
}
