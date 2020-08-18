package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Decodes the body in a string representing the json values. (Used to calculate the hash)
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

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Can't marshall request body.")
		return "", err
	}

	result := string(jsonBody)
	fmt.Printf("ReqBody:%v\n", result)

	return result, nil

}
