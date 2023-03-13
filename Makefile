.PHONY: all

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP_EXECUTABLE="raccoon"
COVER_FILE="/tmp/coverage.out"

all: setup compile

# Setups
setup: copy-config
	make update-deps

update-deps:
	go mod tidy -v
	go mod vendor

copy-config:
	cp .env.sample .env

# Build Lifecycle
compile:
	go build -o $(APP_EXECUTABLE)

build: copy-config update-deps compile

install:
	go install $(ALL_PACKAGES)

start: build
	./$(APP_EXECUTABLE)

# Utility

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

# Tests

test: lint
	ENVIRONMENT=test go test $(shell go list ./... | grep -v "vendor" | grep -v "integration") -v
	@go list ./... | grep -v "vendor" | grep -v "integration" | xargs go test -count 1 -cover -short -race -timeout 1m -coverprofile ${COVER_FILE}
	@go tool cover -func ${COVER_FILE} | tail -1 | xargs echo test coverage:

test-bench: # run benchmark tests
	@go test $(shell go list ./... | grep -v "vendor") -v -bench ./... -run=^Benchmark

test_ci: setup test

# Docker Run

docker-run:
	docker-compose build
	docker-compose up -d

docker-stop:
	docker-compose stop

docker-start:
	docker-compose start