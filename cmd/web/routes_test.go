package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {

	req := routes()

	switch value := req.(type) {
	case http.Handler: // do nothing
	default:
		t.Error(fmt.Sprintf("expected http.Handler, but got %T", value))
	}
}
