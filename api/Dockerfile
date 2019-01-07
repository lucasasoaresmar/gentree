FROM golang:1.11.2

ADD ./src /go/src/app
WORKDIR /go/src/app

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

RUN go install app

ENTRYPOINT /go/bin/app

EXPOSE 3000