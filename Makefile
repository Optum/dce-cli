all: update_mocks test build

test:
	go test -count=1 -v ./...

test_functional:
	go test -count=1 -v tests/functional/

test_unit:
	go test -count=1 -v tests/unit/

.PHONY: mocks
mocks:
	rm -rf mocks/*
	mockery -all

build:
	go build .