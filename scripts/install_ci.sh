#!/usr/bin/env bash
set -euxo pipefail

export GOBIN=$(dirname `which go`)

#export GO111MODULE=off

go get github.com/jstemmer/go-junit-report
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml
go get github.com/matm/gocov-html
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.24.0
go get -u github.com/securego/gosec/cmd/gosec
