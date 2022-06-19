package handler

import (
	"net/http"
	"time"
)

type CustomHandler struct{}

func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC822)
	w.Write([]byte("The time is: " + tm))
}
