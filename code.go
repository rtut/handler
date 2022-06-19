package handler

import (
	"net/http"
	"time"
)

type Handler struct{}

func (ch *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC822)
	w.Write([]byte("The time is: " + tm))
}
gi