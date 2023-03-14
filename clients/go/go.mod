module github.com/goto/raccoon/clients/go

go 1.16

require (
	github.com/stretchr/testify v1.8.2
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.29.0
	honnef.co/go/tools v0.4.2
)

require (
	buf.build/gen/go/gotocompany/proton/grpc/go v1.3.0-20230313110213-9a3d240d5293.1
	buf.build/gen/go/gotocompany/proton/protocolbuffers/go v1.29.0-20230313110213-9a3d240d5293.1
	github.com/gojek/heimdall/v7 v7.0.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
)
