FROM golang:1.21-alpine AS builder

LABEL maintainer="sunist-c"

WORKDIR /app

COPY . /app

RUN apk add --no-cache make && go mod tidy && make build_linux

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/ceobebot-qqchanel_linux-amd64 /app/ceobebot

COPY ./data /app/data

RUN chmod +x /app/ceobebot

CMD ["./ceobebot"]