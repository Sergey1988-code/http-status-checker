BINARY_NAME=http-status-checker

build:
	go build -o $(BINARY_NAME) ./cmd/http-status-checker

build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux ./cmd/http-status-checker
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows.exe ./cmd/http-status-checker
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos ./cmd/http-status-checker

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)*