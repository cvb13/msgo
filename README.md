# msgo - (Mock Server GO)

# Run with GO:
- go run app/main.go
- listening at port 3000

# Run with Docker
- docker build . -t msgo-docker
- docker run -p 3000:3000 msgo-docker

# Postman
- Postman collection included at ***/postman*** folder for testing

# Add mock request
- Method: POST
- URL: http://{{host}}:3000/mocks/add
- Request Body:
```
{
    "url" : "\/test", 
    "requestBody":"{\key\":\"value\"}",
    "requestMethod" : "GET",
    "requestHeaders":{},
    "responseBody":"{\"result\":\"success\"}",
    "responseCode":200,
    "responseHeaders":{"Content-Type":"application/json"},
    "override":false
}
```

Succesful response: 200 OK

# Add multiple mock requests
- Method: POST
- URL: http://{{host}}:3000/mocks/addAll
- Request Body:
```
[
    {
        "url": "/get",
        "requestBody": "",
        "requestMethod": "GET",
        "requestHeaders": {},
        "responseBody": "{\"result\":\"get-success\"}",
        "responseCode": 200,
        "responseHeaders": {
            "Content-Type": "application/json"
        },
        "override": false
    },
    {
        "url": "/put",
        "requestBody": "{\"input\":1}",
        "requestMethod": "PUT",
        "requestHeaders": {},
        "responseBody": "{\"result\":\"put-success\"}",
        "responseCode": 201,
        "responseHeaders": {
            "Content-Type": "application/json"
        },
        "override": false
    }
]
```

Succesful response: 200 OK

# Export in-memory mocks
- Method: GET
- URL: http://{{host}}:3000/mocks/export?fileName=fileName.json
- Response:
```
[
    {
        "url": "/get",
        "requestBody": "",
        "requestMethod": "GET",
        "requestHeaders": {},
        "responseBody": "{\"result\":\"get-success\"}",
        "responseCode": 200,
        "responseHeaders": {
            "Content-Type": "application/json"
        },
        "override": false
    },
    {
        "url": "/put",
        "requestBody": "{\"input\":1}",
        "requestMethod": "PUT",
        "requestHeaders": {},
        "responseBody": "{\"result\":\"put-success\"}",
        "responseCode": 201,
        "responseHeaders": {
            "Content-Type": "application/json"
        },
        "override": false
    }
]
```

Succesfull response code: 200 OK

# Dependencies
- This project uses https://github.com/gorilla/mux to build the handlers and routers for the mocked requests.
