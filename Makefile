NAME="github.com/raystack/raccoon"
COMMIT := $(shell git rev-parse --short HEAD)
TAG := "$(shell git rev-list --tags --max-count=1)"
VERSION := "$(shell git describe --tags ${TAG})"
BUILD_DIR=dist
PROTON_COMMIT := "ccbf219312db35a934361ebad895cb40145ca235"

.PHONY: all build clean test tidy vet proto setup format generate

all: clean test build format lint

tidy:
	@echo "Tidy up go.mod..."
	@go mod tidy -v

lint: ## Lint checker
	@echo "Running lint checks using golangci-lint..."
	@golangci-lint run

clean: tidy ## Clean the build artifacts
	@echo "Cleaning up build directories..."
	@rm -rf $coverage.out ${BUILD_DIR}

proto: ## Generate the protobuf files
	@echo "Generating protobuf from raystack/proton"
	@echo " [info] make sure correct version of dependencies are installed using 'make install'"
	@buf generate https://github.com/raystack/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --template buf.gen.yaml --path raystack/raccoon -v
	@cp -R proto/raystack/raccoon/v1beta1/* proto/ && rm -Rf proto/raystack
	@echo "Protobuf compilation finished"

setup: ## Install required dependencies
	@echo "> Installing dependencies..."
	go mod tidy
	go install github.com/bufbuild/buf/cmd/buf@v1.23.0
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

config: ## Generate the sample config file
	@echo "Initializing sample server config..."
	@cp .env.sample .env

build: ## Build the raccoon binary
	@echo "Building raccoon version ${VERSION}..."
	go build -ldflags "-X ${NAME}/config.Version=${VERSION}" ${NAME}
	@echo "Build complete"

install:
	@echo "Installing Raccoon to ${GOBIN}..."
	go install -ldflags "-X ${NAME}/config.Version=${VERSION}" ${NAME}
	@go install

test: ## Run the tests
	go test $(shell go list ./... | grep -v "vendor" | grep -v "integration") -v

test-bench: # run benchmark tests
	@go test $(shell go list ./... | grep -v "vendor") -v -bench ./... -run=^Benchmark ]

vendor: ## Update the vendor directory
	@echo "Updating vendor directory..."
	@go mod vendor

docker-run:
	docker compose build
	docker compose up -d