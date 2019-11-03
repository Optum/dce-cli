all: update_mocks test build

test:
	go test -count=1 -v ./...

test_functional:
	go test -count=1 -v tests/functional/*

update_mocks:
	rm -rf mocks/*
	mockery -all

build:
	go build .