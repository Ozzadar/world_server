FROM golang

MAINTAINER Paul Mauviel

ADD . /go/src/github.com/ozzadar/world_server
RUN go install github.com/ozzadar/world_server
ENTRYPOINT /go/bin/world_server

EXPOSE 1337