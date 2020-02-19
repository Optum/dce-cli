export GOBIN=$(dirname `which go`)

export GO111MODULE=off

go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
go get -u github.com/securego/gosec/cmd/gosec