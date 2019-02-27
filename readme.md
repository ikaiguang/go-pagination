# go-pagination

The fantastic pagination library for Golang.

> demo [go-pagination-example](https://github.com/ikaiguang/go-pagination-example) 

[![Go Report Card](https://goreportcard.com/badge/github.com/ikaiguang/go-pagination)](https://goreportcard.com/report/github.com/ikaiguang/go-pagination)
[![GoDoc](https://godoc.org/github.com/ikaiguang/go-pagination?status.svg)](https://godoc.org/github.com/ikaiguang/go-pagination)

## Install

install && get source code

`go get -v github.com/ikaiguang/go-pagination`

> go_1.11+ : `GO111MODULE=off go get -v github.com/ikaiguang/go-pagination`

## Test

before run test : you must rewrite database connection and generate test data

```

go get -v github.com/ikaiguang/go-pagination-example
# GO111MODULE=off go get -v github.com/ikaiguang/go-pagination-example

go test -v $GOPATH/src/github.com/ikaiguang/go-pagination-example

```

## Overview

see [pagination.md](pagination.md)

## License

Released under the [MIT License](License)