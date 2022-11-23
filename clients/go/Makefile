.PHONY: all
all : setup fmt staticcheck test

setup:
	@echo "installing the dependencies..."
	go get go.buf.build/odpf/gw/odpf/proton
	go get golang.org/x/lint
	go get google.golang.org/protobuf
	go get github.com/gorilla/websocket
	go get github.com/stretchr/testify
	go get github.com/stretchr/testify
	go get honnef.co/go/tools

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