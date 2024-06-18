NAME="github.com/raystack/raccoon"
VERSION := "$(shell git describe --tags ${TAG})-next"
PROTON_COMMIT := "a4240deecb8345e0e95261f22288f937422594b7"

.PHONY: all build clean test tidy vet proto setup format

all: clean test lint build

tidy:
	@echo "Tidy up go.mod..."
	@go mod tidy -v

lint: ## Lint checker
	@echo "Running lint checks using golangci-lint..."
	@golangci-lint run

clean: tidy ## Clean the build artifacts
	@echo "Cleaning up build directories..."
	@rm -rf $coverage.out raccoon

proto: ## Generate the protobuf files
	@echo "Generating protobuf from raystack/proton"
	@buf generate https://github.com/raystack/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --path raystack/raccoon -v
	@echo "Protobuf compilation finished"

setup: ## Install required dependencies
	@echo "> Installing dependencies..."
	go mod tidy
	go install github.com/bufbuild/buf/cmd/buf@v1.33.0

config: ## Generate the sample config file
	@echo "Initializing sample server config..."
	@cp .env.sample .env

build: ## Build the raccoon binary
	@echo "Building racccoon version ${VERSION}..."
	go build 
	@echo "Build complete"

install:
	@echo "Installing Raccoon to ${GOBIN}..."
	@go install

test: ## Run the tests
	go test $(shell go list ./... | grep -v "vendor" | grep -v "integration") -v

test-bench: # run benchmark tests
	@go test $(shell go list ./... | grep -v "vendor") -v -bench ./... -run=^Benchmark ]

docker-run:
	docker compose build
	docker compose up -d