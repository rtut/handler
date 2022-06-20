package handler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/httptest"
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
	gock.DisableNetworking()
	gock.New("http://www.foo.com").
		MatchHeader("Content-Length", "545").
		Get("/pic1").
		Persist().
		Reply(200).
		JSON(nil)

	body, _ := json.Marshal("http://www.foo.com/pic1\nhttp://www.foo.com/pic1\n")
	request, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	chs.Suite.NoError(err)
	chs.handler.ServeHTTP(res, request)
	gock.EnableNetworking()
	gock.Off()
}

func (chs *customHandlerSuit) TestTimeout() {

}

func (chs *customHandlerSuit) TestLimitURLs() {

}

func (chs *customHandlerSuit) TestWrongMethod() {

}
