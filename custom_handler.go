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

func parseURLs(w http.ResponseWriter, r *http.Request) ([][]byte, error) {
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't parse request body", http.StatusBadRequest)
		return [][]byte{}, err
	}
	responseData = bytes.TrimRight(responseData, "\n")
	count := bytes.Count(responseData, []byte("\n"))
	if count > limitRequest {
		msg := "number of you urls over 100"
		http.Error(w, msg, http.StatusBadRequest)
		return nil, fmt.Errorf(msg)
	}
	urls := bytes.Split(responseData, []byte("\n"))
	return urls, nil
}
