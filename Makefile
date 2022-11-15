.PHONY: all

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP_EXECUTABLE="out/raccoon"
COVER_FILE="/tmp/coverage.out"

all: install-protoc setup compile

# Setups
setup: generate-proto copy-config
	make update-deps

install-protoc:
	@echo "> installing dependencies"
	go get -u github.com/golang/protobuf/proto@v1.4.3
	go get -u github.com/golang/protobuf/protoc-gen-go@v1.4.3
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

update-deps:
	go mod tidy -v
	go mod vendor

copy-config:
	cp .env.sample .env

PROTO_PACKAGE=/proto
generate-proto:
	rm -rf .temp
	mkdir -p .temp
	curl -o .temp/proton.tar.gz -L http://api.github.com/repos/odpf/proton/tarball/main; tar xvf .temp/proton.tar.gz -C .temp/ --strip-components 1
	protoc --proto_path=.temp/ .temp/odpf/raccoon/v1beta1/raccoon.proto --go_out=./ --go_opt=paths=import --go_opt=Modpf/raccoon/v1beta1/raccoon.proto=$(PROTO_PACKAGE)
	protoc --proto_path=.temp/ .temp/odpf/raccoon/v1beta1/raccoon.proto  --go-grpc_opt=paths=import --go-grpc_opt=Modpf/raccoon/v1beta1/raccoon.proto=$(PROTO_PACKAGE) --go-grpc_out=./

# Build Lifecycle
compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

build: copy-config update-deps compile

install:
	go install $(ALL_PACKAGES)

start: build
	./$(APP_EXECUTABLE)

clean: ## Clean the builds
	rm -rf out/

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
	@go test -bench ./...

test_ci: install-protoc setup test

# Docker Run

docker-run:
	docker-compose build
	docker-compose up -d

docker-stop:
	docker-compose stop

docker-start:
	docker-compose start
