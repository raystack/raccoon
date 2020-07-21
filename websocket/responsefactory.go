package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"

	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
)

func createSuccessResponse(requestGuid string) de.EventResponse {
	response := de.EventResponse{
		Status:   de.Status_SUCCESS,
		Code:     de.Code_OK,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data: map[string]string{
			"req_guid": requestGuid,
		},
	}
	return response
}

func createBadrequestResponse(err error) []byte {
	response := de.EventResponse{
		Status:   de.Status_ERROR,
		Code:     de.Code_BAD_REQUEST,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot deserialize request: %s", err),
		Data:     nil,
	}
	badrequestResp, _ := proto.Marshal(&response)
	return badrequestResp
}

func createEmptyErrorResponse(errCode de.Code) []byte {
	resp := de.EventResponse{
		Status:   de.Status_ERROR,
		Code:     errCode,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data:     nil,
	}
	duplicateConnResp, _ := proto.Marshal(&resp)
	return duplicateConnResp
}
