FROM golang:1.12

WORKDIR /go/src/github.com/mtlynch/whatgotdone

RUN  go get github.com/codegangsta/gin

ENTRYPOINT gin --port 3000 --appPort 3001 --all run main.go