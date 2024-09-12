package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"testing/iotest"

	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/core/collector"
	"github.com/raystack/raccoon/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/raystack/raccoon/proto"
)

func TestMain(m *testing.M) {
	logger.SetOutput(io.Discard)
	os.Exit(m.Run())
}

type apiResponse struct {
	Status pb.Status `json:"status"`
	Code   pb.Code   `json:"code"`
	Reason string    `json:"reason"`
}

func TestHandler(t *testing.T) {

	var testCases = []struct {
		Desc      string
		Req       func() *http.Request
		Collector func() collector.Collector
		Response  *apiResponse
		Status    int
		AckType   config.AckType
	}{
		{
			Desc: "should return an error if reading request body fails",
			Req: func() *http.Request {
				e := fmt.Errorf("simulated error")
				rr := httptest.NewRequest("POST", "/api/v1/events", iotest.ErrReader(e))
				rr.Header.Set("Content-Type", "application/json")
				return rr
			},
			Collector: func() collector.Collector { return nil },
			Response: &apiResponse{
				Code:   pb.Code_CODE_INTERNAL_ERROR,
				Status: pb.Status_STATUS_ERROR,
				Reason: "deserialization failure",
			},
			Status: http.StatusInternalServerError,
		},
		{
			Desc: "should return an error if request body is malformed",
			Req: func() *http.Request {
				payload := "}{}"
				rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
				rr.Header.Set("Content-Type", "application/json")
				return rr
			},
			Collector: func() collector.Collector { return nil },
			Response: &apiResponse{
				Code:   pb.Code_CODE_BAD_REQUEST,
				Status: pb.Status_STATUS_ERROR,
				Reason: "deserialization failure",
			},
			Status: http.StatusBadRequest,
		},
		{
			Desc: "should return an error if content-type is unrecognised",
			Req: func() *http.Request {
				payload := "}{}"
				rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
				return rr
			},
			Collector: func() collector.Collector { return nil },
			Response: &apiResponse{
				Code:   pb.Code_CODE_BAD_REQUEST,
				Status: pb.Status_STATUS_ERROR,
				Reason: "invalid content type",
			},
			Status: http.StatusBadRequest,
		},
		{
			Desc: "should return an error if collector fails to consume request (ack type = sync)",
			Req: func() *http.Request {
				payload := "{}"
				rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
				rr.Header.Set("Content-Type", "application/json")
				return rr
			},
			Collector: func() collector.Collector {
				mockCollector := &collector.MockCollector{}
				mockCollector.On("Collect", mock.Anything, mock.Anything).
					Return(nil).
					Once().
					Run(func(args mock.Arguments) {
						args.Get(1).(*collector.CollectRequest).AckFunc(fmt.Errorf("simulated error"))
					})
				return mockCollector
			},
			Response: &apiResponse{
				Code:   pb.Code_CODE_INTERNAL_ERROR,
				Status: pb.Status_STATUS_ERROR,
				Reason: "cannot publish events: simulated error",
			},
			Status:  http.StatusInternalServerError,
			AckType: config.AckTypeSync,
		},
		{
			Desc: "should successfully process event sent (ack type = sync)",
			Req: func() *http.Request {
				payload := "{}"
				rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
				rr.Header.Set("Content-Type", "application/json")
				return rr
			},
			Collector: func() collector.Collector {
				mockCollector := &collector.MockCollector{}
				mockCollector.On("Collect", mock.Anything, mock.Anything).
					Return(nil).
					Once().
					Run(func(args mock.Arguments) {
						args.Get(1).(*collector.CollectRequest).AckFunc(nil)
					})
				return mockCollector
			},
			Response: &apiResponse{
				Code:   pb.Code_CODE_OK,
				Status: pb.Status_STATUS_SUCCESS,
			},
			Status:  http.StatusOK,
			AckType: config.AckTypeSync,
		},
		{
			Desc: "should successfully process event sent (ack type = async)",
			Req: func() *http.Request {
				payload := "{}"
				rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
				rr.Header.Set("Content-Type", "application/json")
				return rr
			},
			Collector: func() collector.Collector {
				mockCollector := &collector.MockCollector{}
				mockCollector.On("Collect", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				return mockCollector
			},
			Response: &apiResponse{
				Code:   pb.Code_CODE_OK,
				Status: pb.Status_STATUS_SUCCESS,
			},
			Status: http.StatusOK,
		},
	}

	for _, testCase := range testCases {

		rw := httptest.NewRecorder()
		h := NewHandler(testCase.Collector())
		h.ackType = testCase.AckType
		h.RESTAPIHandler(rw, testCase.Req())

		assert.Equal(t, testCase.Status, rw.Code)

		res := &apiResponse{}
		assert.Nil(t, json.NewDecoder(rw.Body).Decode(res))
		assert.Equal(t, testCase.Response, res)
	}
}
