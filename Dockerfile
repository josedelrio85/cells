FROM golang:alpine AS builder

RUN apk update && apk add --no-cache \
git \
ca-certificates \
&& update-ca-certificates

COPY . $GOPATH/src/github.com/bysidecar/leads/
WORKDIR $GOPATH/src/github.com/bysidecar/leads/


RUN go get -d -v
RUN go build -o /go/bin/leads


FROM alpine
COPY --from=builder /go/bin/leads /go/bin/leads

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/go/bin/leads"]
