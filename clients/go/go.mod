module github.com/raystack/raccoon/clients/go

go 1.21

toolchain go1.22.4

require (
	github.com/stretchr/testify v1.8.4
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	honnef.co/go/tools v0.3.3
)

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gojek/valkyrie v0.0.0-20180215180059-6aee720afcdf // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20220218215828-6cf2b201936e // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	buf.build/gen/go/raystack/proton/grpc/go v1.4.0-20240713100241-5efa7d29c01b.2
	buf.build/gen/go/raystack/proton/protocolbuffers/go v1.34.2-20240713100241-5efa7d29c01b.2
	github.com/gojek/heimdall/v7 v7.0.3
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.0
)
