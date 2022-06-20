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
		w.Write([]byte("blabla\n"))
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
		fmt.Println(err)
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
		w.Write([]byte("problemmms"))
		return [][]byte{}, fmt.Errorf("")
	}
	urls := bytes.Split(responseData, []byte("\n"))
	return urls, nil
}

