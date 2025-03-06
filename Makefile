run:
	go run main.go

build:
	docker build -t mock-server .

run-docker:
	docker run -p 8080:8080 mock-server

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-short:
	go test -v -short ./...

lint:
	golangci-lint run

clean:
	rm -f coverage.out coverage.html

.PHONY: run build run-docker test test-coverage test-short lint clean

