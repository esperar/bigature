package main

import (
	"bitature/controller/news"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndSearch(t *testing.T) {
	var res *httptest.ResponseRecorder
	var req *http.Request
	var assert = assert.New(t)
	var r = mux.NewRouter()

	news.NewController(r)

	// 생성
	res = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/news", strings.NewReader(`{"title": "Hello", "author": "JehwanYoo", "Content": "This is Test"}`))

	r.ServeHTTP(res, req)
	assert.Equal(http.StatusOK, res.Code)

	// 조회
	res = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/news", nil)
	r.ServeHTTP(res, req)

	var response struct {
		Data []struct {
			Id      string
			Title   string
			Author  string
			Content string
		}
		Status int
		Error  interface{}
	}

	err := json.NewDecoder(res.Body).Decode(&response)

	assert.Nil(err)
	assert.Equal(http.StatusOK, res.Code)
	assert.Nil(response.Error)
	assert.Equal(http.StatusOK, response.Status)
	assert.Equal("Hello", response.Data[0].Title)
	assert.Equal("JehwanYoo", response.Data[0].Author)
	assert.Equal("This is Test", response.Data[0].Content)
}
