package handler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct{}

func (ch *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("blabla\n"))

	}
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	count := bytes.Count(responseData, []byte{'\n'})
	if count > 14 {
		w.Write([]byte("problemmms"))
	}
}

type URL struct {
	url        string
	sizeAnswer uint
}

func (ch *Handler) readBody() {

}
