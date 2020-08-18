package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var mocks []RequestMock

func main() {
	router.HandleFunc("/mocks/add", AddMockHandler).Methods("POST")
	router.HandleFunc("/mocks/addAll", AddAllMockHandler).Methods("POST")
	router.HandleFunc("/mocks/getAll", GetAllMockHandler).Methods("GET")
	router.HandleFunc("/mocks/export", ExportMockHandler).Queries("fileName", "{fileName}").Methods("GET")
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

	addSingleMockRequest(newMock, w)
	w.WriteHeader(201)
}

func AddAllMockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newMocks []RequestMock
	err := json.NewDecoder(r.Body).Decode(&newMocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range newMocks {
		addSingleMockRequest(newMocks[i], w)
	}
	w.WriteHeader(201)
}

//Handler to dynamic handling new requests
func DynamicMockHandler(w http.ResponseWriter, r *http.Request) {

	reqBody, err := decodeBody(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mock := RequestMock{URL: r.URL.Path, RequestMethod: r.Method, RequestBody: reqBody}
	fmt.Printf("Mock:%+v\n", &mock)
	fmt.Printf("ExpectedHash:%x\n", mock.hash())

	if !exists(mock) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	mock = getRequestMock(mock.Hash)
	fmt.Printf("ResponseBody:%s\n", mock.ResponseBody)

	for key, val := range mock.ResponseHeaders {
		w.Header().Set(key, val)
	}
	w.WriteHeader(mock.ResponseCode)
	w.Write([]byte(mock.ResponseBody))
}

//Handler to dynamic handling new requests
func GetAllMockHandler(w http.ResponseWriter, r *http.Request) {
	responseBody, err := json.Marshal(mocks)
	if err != nil {
		http.Error(w, "Error converting mocks to json",
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

//Handler to export current mocks into json file
func ExportMockHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("fileName")
	err := saveMocksToFile(fileName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error saving file.",
			http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, fileName)
}

/** Helper methods **/

// Save mocks to file
func saveMocksToFile(fileName string) error {
	fmt.Printf("FileName: %s\n: ", fileName)
	content, err := json.Marshal(mocks)
	if err != nil {
		fmt.Println("Error marshalling mocks")
		return err
	}
	return ioutil.WriteFile(fileName, content, 0666)
}

//Adds a single request mock into the mocks array
func addSingleMockRequest(newMock RequestMock, w http.ResponseWriter) {
	fmt.Printf("Hash: %x\n", newMock.hash())

	if exists(newMock) {
		if newMock.Override {
			replaceMock(newMock)
		} else {
			http.Error(w, "Duplicated", http.StatusConflict)
			return
		}
	}

	mocks = append(mocks, newMock)

	//handle new path being added
	router.HandleFunc(newMock.URL, DynamicMockHandler).Methods(newMock.RequestMethod)
	fmt.Printf("%+v\n", &newMock)
}

//Validates if the provided mock exists on the mocks slice based on the hash
func exists(mock RequestMock) bool {
	for i := range mocks {
		if mock.Hash == mocks[i].Hash {
			return true
		}
	}
	return false
}

//Retrieves the request mock correspoding to the hash
func getRequestMock(hash [32]byte) RequestMock {
	var result RequestMock
	for i := range mocks {
		if hash == mocks[i].Hash {
			result = mocks[i]
		}
	}
	return result
}

//Replaces mock based on hash
func replaceMock(mock RequestMock) {
	for i := range mocks {
		if mock.Hash == mocks[i].Hash {
			fmt.Printf("ReplacingMock:%+v\n", mock.Hash)
			mocks[i] = mock
			return
		}
	}
}

// Generates hash based on request URL, Method and Body
func (r *RequestMock) hash() [32]byte {
	var buffer bytes.Buffer

	//URL + Method + Body
	buffer.WriteString(r.URL)
	buffer.WriteString(r.RequestMethod)
	buffer.WriteString(r.RequestBody)
	generatedHash := sha256.Sum256([]byte(buffer.String()))

	r.Hash = generatedHash
	return generatedHash
}
