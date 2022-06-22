package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Handler struct{}

func (ch *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "try another method", http.StatusMethodNotAllowed)
		return
	}

	urls, err := parseURLs(w, r)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	res := make(chan int64)

	for _, url := range urls {
		go getContentLength(&wg, res, url)
	}

	go func() {
		wg.Wait()
	}()

	for size := range res {
		r := strconv.FormatInt(size, 10) + "\n"
		w.Write([]byte(r))
	}
	close(res)
}

func getContentLength(wg *sync.WaitGroup, transport chan int64, rawURL []byte) {
	url := string(rawURL)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	transport <- resp.ContentLength
	wg.Done()
}

func parseURLs(w http.ResponseWriter, r *http.Request) ([][]byte, error) {
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return [][]byte{}, err
	}
	responseData = bytes.TrimRight(responseData, "\n")
	count := bytes.Count(responseData, []byte("\n"))
	if count > 14 {
		msg := "number of you urls over 100"
		http.Error(w, msg, http.StatusBadRequest)
		return nil, fmt.Errorf(msg)
	}
	urls := bytes.Split(responseData, []byte("\n"))
	return urls, nil
}
