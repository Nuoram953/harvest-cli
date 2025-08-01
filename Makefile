.PHONY: build run clean test

build:
	go build -o bin/harvest-cli .

run:
	go run .

clean:
	rm -rf bin/

test:
	go test ./...

install:
	go install .
