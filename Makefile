VERSION := $(shell git describe --always --long --dirty)

all: mocks test build

# Generate client code from swagger in a local dce repo (make openapi DCE_REPO=/path/to/dce)
openapi:
	echo "\nMANUAL STEP: Install goswagger cli tool if needed: https://goswagger.io/install.html\n"
	swagger flatten --with-expand $(DCE_REPO)/modules/swagger.yaml >> $(PWD)/out.yaml
	swagger generate client -f $(PWD)/out.yaml --skip-validation -t $(PWD)


# Generate interfaces for OpenApi clients so they can be mocked.
ifaces:
	echo "\nMANUAL STEP: Install interfacer if needed: `go install github.com/rjeczalik/interfaces/cmd/interfacer`\n"
	interfacer -for github.com/Optum/dce-cli/client/operations.Client -as APIer -o internal/util/ifaces.go
	echo "\nMANUAL STEP: Update the package name of internal/util/ifaces.go to package util\n"

.PHONY: mocks
mocks:
	rm -rf mocks/*
	mockery -all

test:
	go test -count=1 -v ./...

cover:
	go test -coverprofile=coverage.out -coverpkg="./pkg/...,./internal/...,./cmd/...,./configs/..."  ./tests/...

test_functional:
	go test -count=1 -v ./tests/functional/

test_unit:
	go test -count=1 -v ./tests/unit/

build:
	go build -ldflags "-X github.com/Optum/dce-cli/cmd.version=${VERSION}" .

.PHONY: docs
docs:
	./update_docs