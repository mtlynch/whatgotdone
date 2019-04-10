FROM golang:1.11.8

EXPOSE $PORT
ENTRYPOINT go run web/main.go