.PHONY: all

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")
APP_EXECUTABLE="out/raccoon"
COVER_FILE="/tmp/coverage.out"

all: install-protoc setup compile

# Setups
setup: generate-proto
	make update-deps

install-protoc:
	@echo "> installing dependencies"
	go get -u github.com/golang/protobuf/proto@v1.4.3
	go get -u github.com/golang/protobuf/protoc-gen-go@v1.4.3

update-deps:
	go mod tidy -v
	go mod vendor

copy-config:
	cp application.yml.sample application.yml

PROTO_PACKAGE=/websocket/proto
generate-proto:
	rm -rf .temp
	mkdir -p .temp
	cd ./.temp; curl -o proton.tar.gz -L http://api.github.com/repos/odpf/proton/tarball/main; tar xvf proton.tar.gz -C . --strip-components 1
	protoc --proto_path=.temp/ $(wildcard .temp/odpf/raccoon/*.proto) --go_out=./ --go_opt=paths=import --go_opt=Modpf/raccoon/Event.proto=$(PROTO_PACKAGE) --go_opt=Modpf/raccoon/EventRequest.proto=$(PROTO_PACKAGE) --go_opt=Modpf/raccoon/EventResponse.proto=$(PROTO_PACKAGE)

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
	ENVIRONMENT=test go test $(shell go list ./... | grep -v "vendor" | grep -v "integration") -p=2 -v
	@go list ./... | grep -v "vendor" | grep -v "integration" | xargs go test -count 1 -cover -short -race -timeout 1m -coverprofile ${COVER_FILE}
	@go tool cover -func ${COVER_FILE} | tail -1 | xargs echo test coverage:

test_ci: install-protoc setup test

# Docker Run

docker-run:
	docker-compose build
	docker-compose up -d

docker-stop:
	docker-compose stop

docker-start:
	docker-compose start
