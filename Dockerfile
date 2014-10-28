FROM golang:1.3.3

MAINTAINER Tyler Rivera <tyler.rivera@gmail.com>

ADD . /go/src/github.com/tjrivera/jim

RUN go get github.com/carlosdp/twiliogo
RUN go install github.com/tjrivera/jim

CMD /go/bin/jim

EXPOSE 8080
