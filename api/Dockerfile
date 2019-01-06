FROM golang:1.11.2

ADD ./src /go/src/app
WORKDIR /go/src/app

RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2/bson
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/txn

RUN go install app

ENTRYPOINT /go/bin/app

EXPOSE 3000