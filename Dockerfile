FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
ENV GOBIN /go/bin
COPY . /go
WORKDIR /go
RUN go install wordfilter
RUN rm -rf pkg src
ENTRYPOINT /go/bin/wordfilter
EXPOSE 50002
