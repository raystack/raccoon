package websocket

import (
	"fmt"
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

func createUnknownrequestResponse(err error) de.EventResponse {
	response := de.EventResponse{
		Status:   de.Status_ERROR,
		Code:     de.Code_BAD_REQUEST,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot deserialize request: %s", err),
		Data:     nil,
	}
	return response
}
