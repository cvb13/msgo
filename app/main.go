package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

//RequestMock struct and functions
type RequestMock struct {
	Hash            [32]byte
	URL             string        `json:"url"`
	RequestBody     string        `json:"requestBody"`
	RequestMethod   string        `json:"requestMethod"`
	RequestHeaders  []interface{} `json:"requestHeaders"`
	ResponseBody    string        `json:"responseBody"`
	ResponseCode    int           `json:"responseCode"`
	ResponseHeaders []interface{} `json:"responseHeaders"`
	Override        bool          `json:"override"`
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mock.ResponseCode)
	w.Write([]byte(mock.ResponseBody))

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

// Decodes the body in a string representing the json values. Used to calculate the hash
func decodeBody(w http.ResponseWriter, r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	reqBody := map[string]interface{}{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Empty body.")
		} else {
			return "", err
		}
	}

	fmt.Printf("ReqBody:%+v\n", &reqBody)

	var result string
	for key, val := range reqBody {
		switch v := val.(type) {
		default:
			fmt.Printf("\nunexpected type %T", v)
		case int:
			n := fmt.Sprintf("{\"%v\":%v}", key, val)
			result = result + n
		case float64:
			n := fmt.Sprintf("{\"%v\":%v}", key, val)
			result = result + n
		case string:
			n := fmt.Sprintf("{\"%v\":\"%v\"}", key, val)
			result = result + n
		}

	}
	return result, nil

}
