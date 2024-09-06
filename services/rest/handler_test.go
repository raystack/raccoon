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

	"github.com/raystack/raccoon/logger"
	"github.com/stretchr/testify/assert"

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

	t.Run("should return an error if reading request body fails", func(t *testing.T) {
		h := NewHandler(nil)

		e := fmt.Errorf("simulated error")
		rr := httptest.NewRequest("POST", "/api/v1/events", iotest.ErrReader(e))
		rr.Header.Set("Content-Type", "application/json")

		rw := httptest.NewRecorder()

		h.RESTAPIHandler(rw, rr)

		assert.Equal(t, rw.Code, http.StatusInternalServerError)

		res := &apiResponse{}
		assert.Nil(t, json.NewDecoder(rw.Body).Decode(res))
		assert.Equal(t, res.Code, pb.Code_CODE_INTERNAL_ERROR)
		assert.Equal(t, res.Status, pb.Status_STATUS_ERROR)
		assert.Equal(t, res.Reason, "deserialization failure")
	})

	t.Run("should return an error if request body is malformed", func(t *testing.T) {
		h := NewHandler(nil)

		payload := "}{}"
		rr := httptest.NewRequest("POST", "/api/v1/events", bytes.NewBufferString(payload))
		rr.Header.Set("Content-Type", "application/json")

		rw := httptest.NewRecorder()

		h.RESTAPIHandler(rw, rr)

		assert.Equal(t, rw.Code, http.StatusBadRequest)

		res := &apiResponse{}
		assert.Nil(t, json.NewDecoder(rw.Body).Decode(res))
		assert.Equal(t, res.Code, pb.Code_CODE_BAD_REQUEST)
		assert.Equal(t, res.Status, pb.Status_STATUS_ERROR)
		assert.Equal(t, res.Reason, "deserialization failure")

	})
}
