package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	homeHandler(response, request)

	got := response.Body.String()
	want := "homepage"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}
