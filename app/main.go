package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var mocks []RequestMock

func main() {
	router.HandleFunc("/mocks/add", AddMockHandler).Methods("POST")
	http.ListenAndServe(":3000", router)
}

func AddMockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMock RequestMock
	err := json.NewDecoder(r.Body).Decode(&newMock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if exists(newMock) {
		http.Error(w, "Duplicated", http.StatusConflict)
		return
	}

	mocks = append(mocks, newMock)

	//handle new path being added
	router.HandleFunc(newMock.URL, DynamicMockHandler).Methods(newMock.RequestMethod)

	fmt.Printf("%+v\n", &newMock)
}

//Handler to dynamic handling new requests
func DynamicMockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mock RequestMock
	mock.URL = r.URL.Path

	if !exists(mock) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	mock = getRequestMock(mock.URL)
	fmt.Printf("ResponseBody:%s\n", mock.ResponseBody)
	w.Write([]byte(mock.ResponseBody))
}

//RequestMock struct and functions
type RequestMock struct {
	URL           string `json:"url"`
	Body          string `json:"body"`
	RequestMethod string `json:"requestMethod"`
	ResponseBody  string `json:"responseBody"`
}

func exists(mock RequestMock) bool {
	for i := range mocks {
		if mock.URL == mocks[i].URL {
			return true
		}
	}
	return false
}

func getRequestMock(url string) RequestMock {
	var result RequestMock
	for i := range mocks {
		if url == mocks[i].URL {
			result = mocks[i]
		}
	}
	return result
}
