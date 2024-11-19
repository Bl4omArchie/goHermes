# syntax=docker/dockerfile:1
FROM golang:alpine

WORKDIR /usr/local/eprint

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o eprint

EXPOSE 8080

CMD ["./eprint"]