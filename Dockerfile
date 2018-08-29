FROM golang:1.10.3 AS builder

ENV GOARCH=amd64
ENV GOOS=linux
ENV CGO_ENABLED=0

ADD . /go/src/github.com/sosedoff/slacklet
WORKDIR /go/src/github.com/sosedoff/slacklet

RUN go get && \
    go build && \
    mv slacklet /

FROM alpine:3.6

RUN \
  apk update && \
  apk add --no-cache ca-certificates openssl bash && \
  update-ca-certificates
WORKDIR /app
COPY --from=builder /slacklet /app/slacklet

CMD ["./slacklet"]