package handler

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type customHandlerSuit struct {
	suite.Suite
	server  *gock.Request
	handler *Handler
}

func TestCustomHandlerSuite(t *testing.T) {
	suite.Run(t, new(customHandlerSuit))

}

func (chs *customHandlerSuit) SetupSuite() {
	chs.handler = &Handler{}

}

func (chs *customHandlerSuit) TearDownSuite() {
}

func (chs *customHandlerSuit) TestResponseOK() {
	gock.New("http://some.url").
		Head("/pic1").
		Persist().
		Reply(200).
		File("1.jpeg")

	urls := "http://some.url/pic1\nhttp://some.url/pic1\n"

	body := bytes.Buffer{}
	body.WriteString(urls)

	request, err := http.NewRequest(http.MethodPost, "", &body)
	request.Header = http.Header{
		"Content-Type": {"text/plain"},
	}
	res := httptest.NewRecorder()
	chs.Suite.NoError(err)
	chs.handler.ServeHTTP(res, request)

	fi, err := os.Stat("1.jpeg")
	chs.NoError(err)
	size := fi.Size()

	gock.Off()

	chs.Equal(fmt.Sprintf("%d\n%d\n", size, size), res.Body.String())
}

func (chs *customHandlerSuit) TestTimeout() {

}

func (chs *customHandlerSuit) TestLimitURLs() {

}

func (chs *customHandlerSuit) TestWrongMethod() {

}
