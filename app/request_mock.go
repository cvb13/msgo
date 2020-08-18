package main

import (
	"bytes"
	"crypto/sha256"
)

//RequestMock struct and functions
type RequestMock struct {
	Hash            [32]byte          `json:"Hash"`
	URL             string            `json:"url"`
	RequestBody     string            `json:"requestBody"`
	RequestMethod   string            `json:"requestMethod"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	ResponseBody    string            `json:"responseBody"`
	ResponseCode    int               `json:"responseCode"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Override        bool              `json:"override"`
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
