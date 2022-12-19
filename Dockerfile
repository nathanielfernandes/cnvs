# syntax=docker/dockerfile:1
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./bin/canvas ./api

EXPOSE 8080
CMD ["./bin/canvas"]