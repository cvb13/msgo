# msgo - (Mock Server GO)

# Run with GO:
- go run app/main.go
- listening at port 3000

# Run with Docker
- docker build . -t msgo-docker
- docker run -p 3000:3000 msgo-docker

# Postman
- Postman collection included at /postman folder for testing

# Dependencies
- This project uses https://github.com/gorilla/mux to build the handlers and routers for the mocked requests.
