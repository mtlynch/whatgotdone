FROM golang:1.11.8

WORKDIR /app/web

EXPOSE $PORT
ENTRYPOINT go run main.go