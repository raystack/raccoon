package ws

import (
	"errors"

	"net/http"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	raccoon "github.com/goto/raccoon/clients/go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/goto/raccoon/clients/go/log"
	"github.com/goto/raccoon/clients/go/retry"
	"github.com/goto/raccoon/clients/go/serializer"
	"github.com/goto/raccoon/clients/go/wire"
)

// New creates the new websocket client with provided options.
func New(options ...Option) (*WS, error) {
	ws := &WS{
		serialize: serializer.PROTO,
		wire:      &wire.ProtoWire{},
		headers:   http.Header{},
		retryMax:  retry.DefaultRetryMax,
		retryWait: retry.DefaultRetryWait,
		logger:    log.Default(),
		acks:      make(chan *raccoon.Response),
	}

	for _, opt := range options {
		opt(ws)
	}

	if err := ws.connect(); err != nil {
		return nil, err
	}

	ws.readMessage()
	return ws, nil
}

// Send sends the events to the raccoon service and returns the req-guid and error if any.
func (ws *WS) Send(events []*raccoon.Event) (string, error) {
	reqId := uuid.NewString()
	ws.logger.Infof("started request, url: %s, req-id: %s", ws.url, reqId)
	defer ws.logger.Infof("ended request, url: %s, req-id: %s", ws.url, reqId)

	e := []*pb.Event{}
	for _, ev := range events {
		// serialize the bytes based on the config
		b, err := ws.serialize(ev.Data)
		if err != nil {
			ws.logger.Errorf("serialize, url: %s, req-id: %s, %+v", ws.url, reqId, err)
			return reqId, err
		}
		e = append(e, &pb.Event{
			EventBytes: b,
			Type:       ev.Type,
		})
	}

	racReq, err := ws.wire.Marshal(&pb.SendEventRequest{
		ReqGuid:  reqId,
		Events:   e,
		SentTime: timestamppb.Now(),
	})
	if err != nil {
		return reqId, err
	}

	err = retry.Do(ws.retryWait, ws.retryMax, func() error {
		err = ws.writeMessage(racReq)
		if err != nil {
			ws.logger.Errorf("send, url: %s, req-id: %s, %+v", ws.url, reqId, err)
			return err
		}
		return nil
	})

	return reqId, err
}

func (ws *WS) EventAcks() <-chan *raccoon.Response {
	return ws.acks
}

// Close closes the connection by sending a close message.
func (ws *WS) Close() {
	if err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		// close the underlying websocket connection.
		ws.Close()
	}
	close(ws.acks)
}

// connnect creates the new client connection.
func (ws *WS) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(ws.url, ws.headers)
	if err != nil {
		return err
	}
	ws.conn = conn
	return nil
}

func (ws *WS) readMessage() {
	go func() {
		for {
			_, b, err := ws.conn.ReadMessage()
			if err != nil {
				return
			}
			resp := pb.SendEventResponse{}
			if err := ws.wire.Unmarshal(b, &resp); err != nil {
				ws.logger.Errorf("wire:unmarshal, url: %s, content-type: %s, %+v", ws.url, ws.wire.ContentType(), err)
				return
			}

			ws.acks <- &raccoon.Response{
				Status:   int32(resp.Status),
				Code:     int32(resp.Code),
				SentTime: resp.SentTime,
				Data:     resp.Data,
			}
		}
	}()
}

func (ws *WS) writeMessage(msg []byte) error {
	switch ws.wire.(type) {
	case *wire.JsonWire:
		if err := ws.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return err
		}
	case *wire.ProtoWire:
		if err := ws.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			return err
		}
	default:
		return errors.New("unsupported wire format")
	}

	return nil
}
