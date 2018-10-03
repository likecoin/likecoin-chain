FROM golang:1.10.3-alpine3.8

ENV DEP_VERSION 0.5.0

# Install Git
RUN apk update && apk upgrade && apk add --no-cache build-base git

# Install dep
ADD https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 /usr/local/bin/dep
RUN chmod +x /usr/local/bin/dep
