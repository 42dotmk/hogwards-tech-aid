package handlers

import (
	"net/http"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}
