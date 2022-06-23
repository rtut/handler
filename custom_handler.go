package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const (
	limitRequest = 100

	errorMethodNotAllowed = "try POST method"
	errorLimitURL         = "number of you urls over %d"
	errorCantParseBody    = "can't parse request body"
)

type Handler struct{}

func (ch *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errorMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	urls, err := parseURLs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	res := make(chan int64)

	for _, url := range urls {
		wg.Add(1)
		go getContentLength(&wg, res, url)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for size := range res {
		r := strconv.FormatInt(size, 10) + "\n"
		w.Write([]byte(r))
	}
}

func getContentLength(wg *sync.WaitGroup, transport chan int64, rawURL []byte) {
	url := string(rawURL)
	resp, err := http.Head(url)
	if err != nil {
		return
	}

	transport <- resp.ContentLength
	wg.Done()
}

func parseURLs(r *http.Request) ([][]byte, error) {
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	responseData = bytes.TrimRight(responseData, "\n")
	count := bytes.Count(responseData, []byte("\n")) + 1
	if count > limitRequest {
		return nil, fmt.Errorf(errorLimitURL, limitRequest)
	}
	urls := bytes.Split(responseData, []byte("\n"))
	return urls, nil
}
