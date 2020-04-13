FROM golang:1.14.1-alpine3.11

ENV GO111MODULE=on

RUN apk update && apk upgrade && \
    apk add --no-cache git openssh gcc g++
