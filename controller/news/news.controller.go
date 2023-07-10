package news

import (
	"bitature/service/news"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CommonResponse struct {
	Data   interface{} `json:data`
	Status int         `json:status`
	Error  interface{} `json:error`
}

func Response(w http.ResponseWriter, data interface{}, status int, err error) {
	var res CommonResponse

	if status == http.StatusOK {
		res.Data = data
		res.Status = status
	} else {
		res.Status = status
		res.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func NewController(router *mux.Router) error {
	err := news.Service.InitService()

	if err != nil {
		return err
	}

	router.HandleFunc("/news/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		raw, err := news.Service.GetOneNews(id)

		if err != nil {
			switch err.Error() {
			case "NOT FOUND":
				Response(w, nil, http.StatusNotFound, errors.New("해당 뉴스가 없습니다."))
			default:
				Response(w, nil, http.StatusInternalServerError, err)
			}
			return
		}

		Response(w, raw, http.StatusOK, nil)

	}).Methods("GET")

	router.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		raws, err := news.Service.GetAllNews()

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
			return
		}

		Response(w, raws, http.StatusOK, nil)
	}).Methods("GET")

	router.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Title   string
			Author  string
			Content string
		}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
		}

		if body.Title == "" || body.Author == "" || body.Content == "" {
			Response(w, nil, http.StatusBadRequest, errors.New("파라미터가 누락되었습니다."))
			return
		}

		err = news.Service.CreateNews(body)

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
			return
		}

		Response(w, "OK", http.StatusOK, nil)

	}).Methods("POST")

	router.HandleFunc("/news/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var body struct {
			Title   string
			Author  string
			Content string
		}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			Response(w, nil, http.StatusInternalServerError, err)
		}

		err = news.Service.UpdateNews(id, body)

		if err != nil {
			switch err.Error() {
			case "NOT FOUND":
				Response(w, nil, http.StatusNotFound, errors.New("해당 뉴스가 없습니다."))
			default:
				Response(w, nil, http.StatusInternalServerError, err)
			}
			return
		}

		Response(w, "OK", http.StatusOK, nil)

	}).Methods("PATCH")

	router.HandleFunc("/news/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err = news.Service.DeleteNews(id)

		if err != nil {
			switch err.Error() {
			case "NOT FOUND":
				Response(w, nil, http.StatusNotFound, errors.New("해당 뉴스가 없습니다."))
			default:
				Response(w, nil, http.StatusInternalServerError, err)
			}
			return
		}

		Response(w, "OK", http.StatusOK, nil)

	}).Methods("DELETE")

	return nil
}
