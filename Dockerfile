FROM golang:wheezy
ADD . /go/src/github.com/dmajere/kogia
WORKDIR /go/src/github.com/dmajere/kogia
ENV GOOS linux
ENV GOARCH amd64
RUN go get
ENTRYPOINT ["/go/src/github.com/dmajere/kogia/make.sh"]
