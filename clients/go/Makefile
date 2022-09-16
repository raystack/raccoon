.PHONY: all
all : install-protoc generate-proto

fmt:
	@go mod tidy
	@go fmt ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck ./...

install-protoc:
	@echo "> installing dependencies..."
	go get -u github.com/golang/protobuf/proto@v1.4.3
	go get -u github.com/golang/protobuf/protoc-gen-go@v1.4.3
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

PROTO_PACKAGE=/proto
generate-proto:
	rm -rf .temp
	mkdir -p .temp
	curl -o .temp/proton.tar.gz -L http://api.github.com/repos/odpf/proton/tarball/main; tar xvf .temp/proton.tar.gz -C .temp/ --strip-components 1
	protoc --proto_path=.temp/ .temp/odpf/raccoon/v1beta1/raccoon.proto --go_out=./ --go_opt=paths=import --go_opt=Modpf/raccoon/v1beta1/raccoon.proto=$(PROTO_PACKAGE)
	protoc --proto_path=.temp/ .temp/odpf/raccoon/v1beta1/raccoon.proto  --go-grpc_opt=paths=import --go-grpc_opt=Modpf/raccoon/v1beta1/raccoon.proto=$(PROTO_PACKAGE) --go-grpc_out=./
	rm -rf .temp