package forms

import (
	"net/url"
	"testing"
)

func TestForms_RequiredField(t *testing.T) {
	nullValue := url.Values{}
	nullValue.Add("first_name", "")
	nullValue.Add("last_name", "")

	form := New(nullValue)
	form.RequiredField("first_name", "last_name")
	if form.IsValid() {
		t.Errorf("fields are blank")
	}

	filled := url.Values{}

	filled.Add("first_name", "allan")
	filled.Add("last_name", "fuentes")

	form = New(filled)
	form.RequiredField("first_name", "last_name")

	if !form.IsValid() {
		t.Error("fields are filled with data")
	}
}

func TestForms_FieldLength(t *testing.T) {
	nullValue := url.Values{}
	nullValue.Add("first_name", "")
	nullValue.Add("last_name", "")

	form := New(nullValue)
	form.FieldLength("first_name", 3)
	if form.IsValid() {
		t.Errorf("expect error cause length is not valid")
	}

	filled := url.Values{}

	filled.Add("first_name", "allan")
	filled.Add("last_name", "fuentes")

	form = New(filled)
	form.FieldLength("first_name", 3)
	if !form.IsValid() {
		t.Error("expected no error, cause value is longer")
	}

	form.Errors.GetErr("first_name")
	form.Errors.GetErr("last_name")
}

func TestForms_IsEmail(t *testing.T) {
	nullValue := url.Values{}
	nullValue.Add("email", "allan@allan")

	form := New(nullValue)
	form.IsEmail("email")
	if form.IsValid() {
		t.Errorf("expect error cause length is not valid")
	}

	filled := url.Values{}

	filled.Add("email", "allan@allan.com")

	form = New(filled)
	form.IsEmail("email")
	if !form.IsValid() {
		t.Error("expected no error, cause value is longer")
	}

	form.Errors.GetErr("email")
}
