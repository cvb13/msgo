package main

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
