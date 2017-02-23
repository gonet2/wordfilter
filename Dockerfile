FROM golang:latest
MAINTAINER xtaci <daniel820313@gmail.com>
COPY . /go/src/wordfilter
RUN go install wordfilter
ENTRYPOINT ["/go/bin/wordfilter"]
EXPOSE 50002
