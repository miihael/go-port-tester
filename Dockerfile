FROM golang:1.20-alpine AS build-env

RUN apk update && apk add --no-cache git

COPY *.go /usr/src/
COPY go.* /usr/src/
WORKDIR /usr/src

ENV CGO_ENABLED=0
RUN go build -ldflags "-extldflags \"-static\"" -o port-tester

FROM alpine:latest
COPY --from=build-env /usr/src/port-tester /usr/local/bin/
COPY docker-entry.sh /usr/local/bin/
ENV SLEEP=86400
ENV PORTS=""
ENV PROTO="tcp"
ENV TARGETS="127.0.0.1"

ENTRYPOINT ["/usr/local/bin/docker-entry.sh"]
