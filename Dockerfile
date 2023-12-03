FROM golang:1.21.4-alpine3.18 AS builder
LABEL maintainer="Joseph Mate"
RUN apk add --no-cache --update git && apk add build-base
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ .
RUN go build -o bin/server


FROM alpine:3.18
RUN apk --no-cache add ca-certificates
# need git for uploading the new file to git
RUN apk add --no-cache git
WORKDIR /app
COPY --from=builder /app/bin/server bin/server
