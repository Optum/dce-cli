all: mocks test build

openapi:
	echo "MANUAL STEP: Install goswagger cli tool if needed: https://goswagger.io/install.html"
	swagger flatten --with-expand $(DCE_REPO)/modules/swagger.yaml >> $(DCE_REPO)/out.yaml
	swagger generate client -f $(DCE_REPO)/out.yaml --skip-validation -t $(PWD)


# Generate interfaces for OpenApi clients so they can be mocked.
ifaces:
	echo "MANUAL STEP: Install interfacer if needed: `go install github.com/rjeczalik/interfaces/cmd/interfacer`"
	interfacer -for github.com/Optum/dce-cli/client/operations.Client -as APIer -o internal/util/ifaces.go
	echo "\nMANUAL STEP: Update the package name of internal/util/ifaces.go to package util"

.PHONY: mocks
mocks:
	rm -rf mocks/*
	mockery -all

test:
	go test -count=1 -v ./...

test_functional:
	go test -count=1 -v ./tests/functional/

test_unit:
	go test -count=1 -v ./tests/unit/

build:
	go build .