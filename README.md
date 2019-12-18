# Leads API

This repository contains an implementation of an API to handle C2C data generated in our LP's and store it.

## Why do we need it

We created this API because we need to rethink the logic implemented in others repositories.

The goal is to create a more maintainable code using the capabilities of Go! language.

## How to run the service

This service has been created with the following GO version:

```bash
go version go1.11.5 darwin/amd64
```

It's a HTTP service that could be run locally on the 5000 port using:

```bash
go run main.go
```

## How to run the tests

```bash
go test ./...
```
