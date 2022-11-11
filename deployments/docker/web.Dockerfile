FROM golang:alpine

# Required because go requires gcc to build
RUN apk add build-base
RUN apk add inotify-tools
RUN apk add git
RUN go install github.com/rubenv/sql-migrate/...@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN echo $GOPATH

COPY . /web
WORKDIR /web

RUN go mod download

CMD sh /web/docker/run.sh