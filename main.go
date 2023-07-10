package main

import (
	"bitature/controller/news"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	err := news.NewController(r)

	if err != nil {
		panic("서버 실행에 실패했습니다.")
	}

	http.ListenAndServe(":3000", r)
}
