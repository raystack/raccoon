package raccoon

import (
	"bytes"
	"io"

	"net/http"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/google/uuid"
)

// NewRest creates the new rest client with provided options.
func NewRest(options ...RestOption) *RestClient {
	rc := &RestClient{
		Wire:       &JsonWire{}, // default
		httpclient: httpclient.NewClient(),
		headers:    http.Header{},
	}

	for _, opt := range options {
		opt(rc)
	}

	return rc
}

// Send sends the events to the raccoon service
func (c *RestClient) Send(events []*Event) (string, *Response, error) {
	reqId := uuid.NewString()

	e := []*pb.Event{}
	for _, ev := range events {
		// serialize the bytes based on the config
		b, err := c.Marshal(ev.Data)
		if err != nil {
			return reqId, nil, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	racReq, err := c.Wire.Marshal(&pb.SendEventRequest{
		ReqGuid: reqId,
		Events:  e,
	})
	if err != nil {
		return reqId, nil, err
	}

	res, err := c.executeRequest(racReq)
	if err != nil {
		return reqId, nil, err
	}

	resp := pb.SendEventResponse{}
	if err := c.Wire.Unmarshal(res, &resp); err != nil {
		return reqId, nil, err
	}

	return reqId, &Response{
		Status:   int32(resp.Status),
		Code:     int32(resp.Code),
		SentTime: resp.SentTime,
		Data:     resp.Data,
	}, nil
}

func (c *RestClient) executeRequest(body []byte) ([]byte, error) {
	c.headers.Set("Content-Type", c.Wire.ContentType())
	resp, err := c.httpclient.Post(c.Url, bytes.NewReader(body), c.headers)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
