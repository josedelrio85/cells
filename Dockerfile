# Builder stage
FROM golang:alpine AS builder

# Install dependencies
RUN apk update && apk add --no-cache \
  git \
  ca-certificates \
  && update-ca-certificates

# Add source files and set the proper work dir
COPY . $GOPATH/src/github.com/josedelrio85/leads/
WORKDIR $GOPATH/src/github.com/josedelrio85/leads/

# Fetch dependencies.
# RUN go get -d -v
# Enable Go Modules
ENV GO111MODULE=on
# Build the binary.
RUN go build -mod=vendor -o /go/bin/leads

# Final image
FROM alpine

# Copy our static executable.
COPY --from=builder /go/bin/leads /go/bin/leads

# Copy the ca-certificates to be able to perform https requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the hello binary.
ENTRYPOINT ["/go/bin/leads"]
