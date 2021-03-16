package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"time"

	pb "raccoon/websocket/proto"
)

func createSuccessResponse(requestGUID string) *pb.EventResponse {
	response := &pb.EventResponse{
		Status:   pb.Status_SUCCESS,
		Code:     pb.Code_OK,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data: map[string]string{
			"req_guid": requestGUID,
		},
	}
	return response
}

func createBadrequestResponse(err error) []byte {
	response := pb.EventResponse{
		Status:   pb.Status_ERROR,
		Code:     pb.Code_BAD_REQUEST,
		SentTime: time.Now().Unix(),
		Reason:   fmt.Sprintf("cannot deserialize request: %s", err),
		Data:     nil,
	}
	badrequestResp, _ := proto.Marshal(&response)
	return badrequestResp
}

func createEmptyErrorResponse(errCode pb.Code) []byte {
	resp := pb.EventResponse{
		Status:   pb.Status_ERROR,
		Code:     errCode,
		SentTime: time.Now().Unix(),
		Reason:   "",
		Data:     nil,
	}
	duplicateConnResp, _ := proto.Marshal(&resp)
	return duplicateConnResp
}
