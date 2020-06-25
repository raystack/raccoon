package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"

	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
)

func createSuccessResponse(request de.EventRequest) de.EventResponse {
	response := de.EventResponse{
		Status:   de.Status_SUCCESS,
		Code:     de.Code_OK,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data: map[string]string{
			"req_guid": request.ReqGuid,
		},
	}
	return response
}

func createBadrequestResponse(err error) de.EventResponse {
	response := de.EventResponse{
		Status:   de.Status_ERROR,
		Code:     de.Code_BAD_REQUEST,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot deserialize request: %s", err),
		Data:     nil,
	}
	return response
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