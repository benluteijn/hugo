FROM golang:1.4-cross

ENV GOPATH /go
ENV USER root

#WORKDIR /go/src/github.com/spf13/hugo
#ADD . /go/src/github.com/spf13/hugo

ENV trigger 1
RUN go get -d -v github.com/spf13/hugo
WORKDIR /go/src/github.com/spf13/hugo

RUN go get -d -v
ADD . /go/src/github.com/spf13/hugo
RUN go build -o hugo main.go

