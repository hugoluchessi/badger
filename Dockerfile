FROM golang:1.12.0-alpine3.9

ENV GO111MODULE=on

RUN apk update && apk upgrade && \
    apk add --no-cache git gcc g++
