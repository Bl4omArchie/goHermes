# syntax=docker/dockerfile:1
FROM golang:alpine

WORKDIR /usr/local/eprint

COPY . .
RUN go mod download
RUN go build -o eprint

EXPOSE 1337

CMD ["./eprint"]