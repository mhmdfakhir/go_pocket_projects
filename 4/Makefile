# @ suppresses the output of the commands
run:
	@go run .

build:
	@go build -o bin/ .

test:
	@go test ./...

cover:
	@go test ./... -cover

lint:
	@golangci-lint run ./...