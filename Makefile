.PHONY: run build launch fmt vet test

run:
	@go run main.go

build: fmt vet test
	@go build -o ./build/tbot ./main.go
	@chmod +x ./build/tbot

launch: build
	@./build/tbot

fmt:
	@go fmt ./...

vet:
	@go vet ./...

test:
	@go test ./...
