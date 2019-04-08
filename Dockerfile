FROM golang:1.11

WORKDIR /go/src/github.com/mtlynch/whatgotdone

RUN go get github.com/codegangsta/gin

EXPOSE 3000

ENTRYPOINT gin --port 3000 --appPort 3001 run main.go