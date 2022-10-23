package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCsrfToken(t *testing.T) {
	var mh myHandler

	req := CsrfToken(mh)

	switch value := req.(type) {
	case http.Handler: // do nothing
	default:
		t.Error(fmt.Sprintf("expected http.Handler, but got %T", value))
	}
}

func TestLoadSession(t *testing.T) {
	var mh myHandler

	req := LoadSession(mh)

	switch value := req.(type) {
	case http.Handler: // do nothing
	default:
		t.Error(fmt.Sprintf("expected http.Handler, but got %T", value))
	}
}
