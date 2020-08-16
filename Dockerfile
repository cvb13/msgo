FROM golang as base

WORKDIR /msgo

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY app/go.mod .
COPY app/go.sum .

RUN go mod download

COPY app/. .

# it will take the flags from the environment
RUN go build

### Certs
FROM alpine:latest as certs
RUN apk --update add ca-certificates

EXPOSE 3000

### App
FROM scratch as app
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base msgo /
ENTRYPOINT ["/msgo"]