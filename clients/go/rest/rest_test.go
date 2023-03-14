package rest

import (
	"io"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"

	raccoon "github.com/goto/raccoon/clients/go"
	"github.com/goto/raccoon/clients/go/serializer"
	"github.com/stretchr/testify/assert"
)

func TestRestClientSend(t *testing.T) {
	assert := assert.New(t)
	var success int32 = 1
	var req_guid string
	clickEvent := `{ "name": "raccoon" }`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("application/json", r.Header.Get("Content-Type"))

		rBody, err := io.ReadAll(r.Body)
		assert.NoError(err)

		req := &pb.SendEventRequest{}
		err = json.Unmarshal(rBody, req)
		assert.NoError(err)

		req_guid = req.ReqGuid

		gotEvent := ""
		json.Unmarshal(req.Events[0].EventBytes, &gotEvent)
		assert.Equal(clickEvent, gotEvent)

		w.WriteHeader(http.StatusOK)
		res := &pb.SendEventResponse{
			Status:   pb.Status_STATUS_SUCCESS,
			Code:     pb.Code_CODE_OK,
			Reason:   "",
			SentTime: time.Now().UTC().Unix(),
			Data: map[string]string{
				"req_guid": req.ReqGuid,
			},
		}
		b, _ := json.Marshal(res)
		w.Write(b)
	}

	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	rc, err := New(
		WithUrl(server.URL),
		WithSerializer(serializer.JSON),
		WithHeader("admin", "admin"))

	assert.NoError(err)

	reqGuid, resp, err := rc.Send([]*raccoon.Event{
		{
			Type: "page",
			Data: `{ "name": "raccoon" }`,
		},
	})

	assert.Equal(req_guid, reqGuid)
	assert.Nil(err)
	assert.Equal(success, resp.Status)
	assert.NotNil(resp.SentTime)
}
