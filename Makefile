.PHONY: build run test clean

build:
	go build -o bin/aggregator cmd/aggregator/main.go

run:
	go run cmd/aggregator/main.go

test:
	go test ./internal/... -v

clean:
	rm -rf bin/
