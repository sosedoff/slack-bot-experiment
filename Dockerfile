FROM golang:1.10.3 AS builder

ADD . /go/src/github.com/sosedoff/slacklet
WORKDIR /go/src/github.com/sosedoff/slacklet
RUN go get && \
    go build && \
    mv slacklet /

FROM alpine:3.6
RUN \
  apk update && \
  apk add --no-cache ca-certificates openssl && \
  update-ca-certificates

COPY --from=builder /slacklet /bin/slacklet

CMD ["/bin/slacklet"]