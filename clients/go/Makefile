.PHONY: all
all : fmt staticcheck test

fmt:
	@echo "code formatting..."
	@go mod tidy
	@go fmt ./...

staticcheck:
	@echo "running static analysis tools..."
	@go run honnef.co/go/tools/cmd/staticcheck ./...

test:
	@echo "running tests..."
	@go test -race -v ./... --cover