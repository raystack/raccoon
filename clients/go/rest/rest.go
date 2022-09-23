package rest

import (
	"bytes"
	"io"

	"net/http"

	pb "go.buf.build/odpf/gw/odpf/proton/odpf/raccoon/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/google/uuid"
	raccoon "github.com/odpf/raccoon/clients/go"
	"github.com/odpf/raccoon/clients/go/serializer"
	"github.com/odpf/raccoon/clients/go/wire"
)

// NewRest creates the new rest client with provided options.
func NewRest(options ...RestOption) (*RestClient, error) {
	rc := &RestClient{
		serialize:  serializer.JSON,
		wire:       &wire.JsonWire{},
		httpclient: httpclient.NewClient(),
		headers:    http.Header{},
	}

	for _, opt := range options {
		opt(rc)
	}

	return rc, nil
}

// Send sends the events to the raccoon service
func (c *RestClient) Send(events []*raccoon.Event) (string, *raccoon.Response, error) {
	reqId := uuid.NewString()

	e := []*pb.Event{}
	for _, ev := range events {
		// serialize the bytes based on the config
		b, err := c.serialize(ev.Data)
		if err != nil {
			return reqId, nil, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	racReq, err := c.wire.Marshal(&pb.SendEventRequest{
		ReqGuid:  reqId,
		Events:   e,
		SentTime: timestamppb.Now(),
	})
	if err != nil {
		return reqId, nil, err
	}

	res, err := c.executeRequest(racReq)
	if err != nil {
		return reqId, nil, err
	}

	resp := pb.SendEventResponse{}
	if err := c.wire.Unmarshal(res, &resp); err != nil {
		return reqId, nil, err
	}

	return reqId, &raccoon.Response{
		Status:   int32(resp.Status),
		Code:     int32(resp.Code),
		SentTime: resp.SentTime,
		Data:     resp.Data,
	}, nil
}

func (c *RestClient) executeRequest(body []byte) ([]byte, error) {
	c.headers.Set("Content-Type", c.wire.ContentType())
	resp, err := c.httpclient.Post(c.url, bytes.NewReader(body), c.headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
