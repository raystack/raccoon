package raccoon

import (
	"io"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"

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
				"req_guid": "test-guid",
			},
		}
		b, _ := json.Marshal(res)
		w.Write(b)
	}

	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	rc := NewRest(
		WithUrl(server.URL),
		WithMarshaler(JSON),
		WithHeader("admin", "admin"))

	reqGuid, resp, err := rc.Send([]*Event{
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
