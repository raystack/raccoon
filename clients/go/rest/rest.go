package rest

import (
	"bytes"
	"fmt"
	"io"

	"net/http"

	pb "go.buf.build/raystack/gw/raystack/proton/raystack/raccoon/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/google/uuid"
	raccoon "github.com/raystack/raccoon/clients/go"
	"github.com/raystack/raccoon/clients/go/log"
	"github.com/raystack/raccoon/clients/go/retry"
	"github.com/raystack/raccoon/clients/go/serializer"
	"github.com/raystack/raccoon/clients/go/wire"
)

// New creates the new rest client with provided options.
func New(options ...Option) (*Rest, error) {
	rc := &Rest{
		serialize:  serializer.JSON,
		wire:       &wire.JsonWire{},
		httpclient: httpclient.NewClient(),
		headers:    http.Header{},
		retryMax:   retry.DefaultRetryMax,
		retryWait:  retry.DefaultRetryWait,
		logger:     log.Default(),
	}

	for _, opt := range options {
		opt(rc)
	}

	return rc, nil
}

// Send sends the events to the raccoon service
func (rc *Rest) Send(events []*raccoon.Event) (string, *raccoon.Response, error) {
	reqId := uuid.NewString()
	rc.logger.Infof("started request, url: %s, req-id: %s", rc.url, reqId)
	defer rc.logger.Infof("ended request, url: %s, req-id: %s", rc.url, reqId)

	e := []*pb.Event{}
	for _, ev := range events {
		// serialize the bytes based on the config
		b, err := rc.serialize(ev.Data)
		if err != nil {
			rc.logger.Errorf("serialize, url: %s, req-id: %s, %+v", rc.url, reqId, err)
			return reqId, nil, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	racReq, err := rc.wire.Marshal(&pb.SendEventRequest{
		ReqGuid:  reqId,
		Events:   e,
		SentTime: timestamppb.Now(),
	})
	if err != nil {
		return reqId, nil, err
	}

	resp := pb.SendEventResponse{}
	err = retry.Do(rc.retryWait, rc.retryMax, func() error {
		b, err := rc.executeRequest(racReq)
		if err != nil {
			return err
		}

		if err := rc.wire.Unmarshal(b, &resp); err != nil {
			rc.logger.Errorf("wire:unmarshal, url: %s, req-id: %s, content-type: %s, %+v", rc.url, reqId, rc.wire.ContentType(), err)
			return err
		}

		if resp.Status != pb.Status_STATUS_SUCCESS {
			return fmt.Errorf("error from raccoon url: %s, req-id: %s, status: %d, code: %d, data: %+v", rc.url, reqId, resp.Status, resp.Code, resp.Data)
		}
		return nil
	})
	if err != nil {
		rc.logger.Errorf("send, url: %s, req-id: %s, %+v", rc.url, reqId, err)
		return reqId, nil, err
	}

	return reqId, &raccoon.Response{
		Status:   int32(resp.Status),
		Code:     int32(resp.Code),
		SentTime: resp.SentTime,
		Data:     resp.Data,
	}, nil
}

func (rc *Rest) executeRequest(body []byte) ([]byte, error) {
	rc.headers.Set("Content-Type", rc.wire.ContentType())
	resp, err := rc.httpclient.Post(rc.url, bytes.NewReader(body), rc.headers)
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
