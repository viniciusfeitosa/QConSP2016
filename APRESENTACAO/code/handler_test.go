package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusOK {
		t.Error("Erro ao retornar o code")
	}
	if w.Body.String() != "Hello, World\n" {
		t.Error("Erro ao retornar o body")
	}
}
