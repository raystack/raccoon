.PHONY: all
all : fmt staticcheck

fmt:
	@go mod tidy
	@go fmt ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck ./...