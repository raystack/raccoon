package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	pbgrpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/raccoon/v1beta1/raccoonv1beta1grpc"
	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var uuid string
var timeout time.Duration
var topicFormat string
var url, wsurl string
var bootstrapServers string
var grpcServerAddr string

func TestMain(m *testing.M) {
	uuid = fmt.Sprintf("%d-test", rand.Int())
	timeout = 20 * time.Second
	topicFormat = os.Getenv("INTEGTEST_TOPIC_FORMAT")
	wsurl = fmt.Sprintf("ws://%v/api/v1/events", os.Getenv("INTEGTEST_HOST"))
	url = fmt.Sprintf("http://%v/api/v1/events", os.Getenv("INTEGTEST_HOST"))
	grpcServerAddr = os.Getenv("GRPC_SERVER_ADDR")
	bootstrapServers = os.Getenv("INTEGTEST_BOOTSTRAP_SERVER")
	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	var err error
	assert.NoError(t, err)
	header := http.Header{
		"X-User-ID": []string{"1234"},
	}
	t.Run("Should response with BadRequest when sending invalid request", func(t *testing.T) {

		wss, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		if err != nil {
			assert.Fail(t, fmt.Sprintf("fail to connect. %v", err))
		}

		defer wss.Close()
		wss.WriteMessage(websocket.BinaryMessage, []byte{1})
		mType, resp, err := wss.ReadMessage()
		r := &pb.SendEventResponse{}
		_ = proto.Unmarshal(resp, r)
		assert.Equal(t, mType, websocket.BinaryMessage)
		assert.Empty(t, err)
		assert.Equal(t, pb.Status_STATUS_ERROR, r.Status)
		assert.Equal(t, pb.Code_CODE_BAD_REQUEST, r.Code)
		assert.NotEmpty(t, r.Reason)
		assert.Equal(t, r.Data, map[string]string{"req_guid": ""})

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should response with BadRequest when sending invalid GRPC request", func(t *testing.T) {
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

		conn, err := grpc.Dial(grpcServerAddr, opts...)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("fail to connect grpc server. %v", err))
		}
		defer conn.Close()

		client := pbgrpc.NewEventServiceClient(conn)
		r, err := client.SendEvent(context.Background(), &pb.SendEventRequest{})

		assert.NotEmpty(t, err)
		assert.Empty(t, r)

	})

	t.Run("Should response with BadRequest when sending invalid json request", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		if err != nil {
			assert.Fail(t, fmt.Sprintf("fail to connect. %v", err))
		}
		defer wss.Close()

		wss.WriteMessage(websocket.TextMessage, []byte(""))

		mType, resp, err := wss.ReadMessage()
		r := &pb.SendEventResponse{}
		_ = json.Unmarshal(resp, r)
		assert.Equal(t, websocket.TextMessage, mType)
		assert.Empty(t, err)
		assert.Equal(t, pb.Status_STATUS_ERROR, r.Status)
		assert.Equal(t, pb.Code_CODE_BAD_REQUEST, r.Code)
		assert.NotEmpty(t, r.Reason)
		assert.Equal(t, r.Data, map[string]string{"req_guid": ""})

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should response with BadRequest when sending HTTP/1.1 request with empty JSON body", func(t *testing.T) {
		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to create http request. %v", err))
			os.Exit(1)
		}
		res, err := client.Do(req)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to connect to http server. %v", err))
			os.Exit(1)
		}
		defer io.Copy(ioutil.Discard, res.Body)
		defer res.Body.Close()
		r := &pb.SendEventResponse{}
		err = json.NewDecoder(res.Body).Decode(r)
		assert.Empty(t, err)
		assert.Equal(t, pb.Status_STATUS_ERROR, r.Status)
		assert.Equal(t, pb.Code_CODE_BAD_REQUEST, r.Code)
		assert.NotEmpty(t, r.Reason)
		assert.Empty(t, r.Data)
	})

	t.Run("Should response with BadRequest when sending HTTP/1.1 request with invalid JSON body", func(t *testing.T) {
		client := &http.Client{}

		bodyBuf := bytes.NewBuffer([]byte{1})
		req, err := http.NewRequest(http.MethodPost, url, bodyBuf)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to create http request. %v", err))
			os.Exit(1)
		}
		res, err := client.Do(req)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to connect to http server. %v", err))
			os.Exit(1)
		}
		defer io.Copy(ioutil.Discard, res.Body)
		defer res.Body.Close()
		r := &pb.SendEventResponse{}
		err = json.NewDecoder(res.Body).Decode(r)
		assert.Empty(t, err)
		assert.Equal(t, pb.Status_STATUS_ERROR, r.Status)
		assert.Equal(t, pb.Code_CODE_BAD_REQUEST, r.Code)
		assert.NotEmpty(t, r.Reason)
		assert.Empty(t, r.Data)
	})

	t.Run("Should response with success when HTTP/1.1 REST JSON request is processed", func(t *testing.T) {
		client := &http.Client{}
		var events []*pb.Event

		eEvent1 := &pb.Event{
			EventBytes: []byte("event_1"),
			Type:       "type_a",
		}
		eEvent2 := &pb.Event{
			EventBytes: []byte("event_2"),
			Type:       "type_b",
		}
		events = append(events, eEvent1)
		events = append(events, eEvent2)
		req := &pb.SendEventRequest{
			ReqGuid:  "1234",
			SentTime: timestamppb.Now(),
			Events:   events,
		}
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(req)

		request, err := http.NewRequest(http.MethodPost, url, buf)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to create http request. %v", err))
			os.Exit(1)
		}
		request.Header.Add("Content-Type", "application/json")
		res, err := client.Do(request)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("failed to connect to http server. %v", err))
			os.Exit(1)
		}
		defer io.Copy(ioutil.Discard, res.Body)
		defer res.Body.Close()
		r := &pb.SendEventResponse{}
		err = json.NewDecoder(res.Body).Decode(r)
		assert.Empty(t, err)
		assert.Equal(t, pb.Code_CODE_OK, r.Code)
		assert.Equal(t, pb.Status_STATUS_SUCCESS, r.Status)
		assert.Equal(t, r.Reason, "")
		assert.Equal(t, r.Data, map[string]string{"req_guid": "1234"})

	})

	t.Run("Should response with success when JSON request is processed", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		if err != nil {
			panic(err)
		}
		defer wss.Close()
		var events []*pb.Event

		eEvent1 := &pb.Event{
			EventBytes: []byte("event_1"),
			Type:       "type_a",
		}
		eEvent2 := &pb.Event{
			EventBytes: []byte("event_2"),
			Type:       "type_b",
		}
		events = append(events, eEvent1)
		events = append(events, eEvent2)
		req := &pb.SendEventRequest{
			ReqGuid:  "1234",
			SentTime: timestamppb.Now(),
			Events:   events,
		}
		bReq, _ := json.Marshal(req)
		wss.WriteMessage(websocket.TextMessage, bReq)

		mType, resp, err := wss.ReadMessage()
		r := &pb.SendEventResponse{}
		_ = json.Unmarshal(resp, r)
		assert.Equal(t, mType, websocket.TextMessage)
		assert.Empty(t, err)
		assert.Equal(t, r.Code.String(), pb.Code_CODE_OK.String())
		assert.Equal(t, r.Status.String(), pb.Status_STATUS_SUCCESS.String())
		assert.Equal(t, r.Reason, "")
		assert.Equal(t, r.Data, map[string]string{"req_guid": "1234"})

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should response with success when correct GRPC request is processed", func(t *testing.T) {
		opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

		conn, err := grpc.Dial(grpcServerAddr, opts...)
		if err != nil {
			assert.Fail(t, fmt.Sprintf("fail to connect grpc server. %v", err))
			os.Exit(1)
		}
		defer conn.Close()

		client := pbgrpc.NewEventServiceClient(conn)
		var events []*pb.Event

		eEvent1 := &pb.Event{
			EventBytes: []byte("event_1"),
			Type:       "type_a",
		}
		eEvent2 := &pb.Event{
			EventBytes: []byte("event_2"),
			Type:       "type_b",
		}
		events = append(events, eEvent1)
		events = append(events, eEvent2)
		req := &pb.SendEventRequest{
			ReqGuid:  "1234",
			SentTime: timestamppb.Now(),
			Events:   events,
		}

		md := metadata.New(map[string]string{"X-User-ID": "1234"})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		r, err := client.SendEvent(ctx, req)
		assert.Empty(t, err)
		assert.Equal(t, r.Code.String(), pb.Code_CODE_OK.String())
		assert.Equal(t, r.Status.String(), pb.Status_STATUS_SUCCESS.String())
		assert.Equal(t, r.Reason, "")
		assert.Equal(t, r.Data, map[string]string{"req_guid": "1234"})

	})

	t.Run("Should response with success when request is processed", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		if err != nil {
			panic(err)
		}

		defer wss.Close()
		var events []*pb.Event

		eEvent1 := &pb.Event{
			EventBytes: []byte("event_1"),
			Type:       "type_a",
		}
		eEvent2 := &pb.Event{
			EventBytes: []byte("event_2"),
			Type:       "type_b",
		}
		events = append(events, eEvent1)
		events = append(events, eEvent2)
		req := &pb.SendEventRequest{
			ReqGuid:  "1234",
			SentTime: timestamppb.Now(),
			Events:   events,
		}
		bReq, _ := proto.Marshal(req)
		wss.WriteMessage(websocket.BinaryMessage, bReq)

		mType, resp, err := wss.ReadMessage()
		r := &pb.SendEventResponse{}
		_ = proto.Unmarshal(resp, r)
		assert.Equal(t, mType, websocket.BinaryMessage)
		assert.Empty(t, err)
		assert.Equal(t, r.Code.String(), pb.Code_CODE_OK.String())
		assert.Equal(t, r.Status.String(), pb.Status_STATUS_SUCCESS.String())
		assert.Equal(t, r.Reason, "")
		assert.Equal(t, r.Data, map[string]string{"req_guid": "1234"})

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should be able to consume published message", func(t *testing.T) {
		t.Run("type_a", func(t *testing.T) {
			t.Parallel()
			c, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": bootstrapServers,
				"group.id":          "my-local-group",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				assert.Fail(t, "setup kafka consumer failed")
			}

			e := c.Subscribe(fmt.Sprintf(topicFormat, "type_a"), nil)
			if e != nil {
				assert.Fail(t, fmt.Sprintf("Pls try again. %v", e))
			}
			timer := time.After(timeout)
			for {
				select {
				case <-timer:
					t.Fatal("timeout")
				default:
					msg, err := c.ReadMessage(timeout)
					if err != nil {
						t.Log("error", err)
						continue
					}
					if string(msg.Value) == "event_1" {
						return
					}
				}
			}
		})

		t.Run("type_b", func(t *testing.T) {
			t.Parallel()
			c, err := kafka.NewConsumer(&kafka.ConfigMap{
				"bootstrap.servers": bootstrapServers,
				"group.id":          "my-local-group",
				"auto.offset.reset": "earliest",
			})

			if err != nil {
				assert.Fail(t, "setup kafka consumer failed")
			}

			e := c.Subscribe(fmt.Sprintf(topicFormat, "type_b"), nil)
			if e != nil {
				assert.Fail(t, fmt.Sprintf("Pls try again. %v", e))
			}
			timer := time.After(timeout)
			for {
				select {
				case <-timer:
					t.Fatal("timeout")
				default:
					msg, err := c.ReadMessage(timeout)
					if err != nil {
						t.Log("error", err)
						continue
					}
					if string(msg.Value) == "event_2" {
						return
					}
				}
			}
		})

	})

	t.Run("Should close connection when client is unresponsive", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		if err != nil {
			assert.Fail(t, err.Error())
		}

		defer wss.Close()
		wss.SetPingHandler(func(appData string) error {
			return nil
		})

		done := make(chan int)
		wss.SetCloseHandler(func(code int, text string) error {
			close(done)
			return nil
		})

		go func() { _, _, err = wss.ReadMessage() }()
		select {
		case <-time.After(timeout):
			break
		case <-done:
			break
		}

		assert.Error(t, err)
	})

	t.Run("Should disconnect subsequence connection from same user when already connected", func(t *testing.T) {
		done := make(chan int)
		_, _, err := websocket.DefaultDialer.Dial(wsurl, header)

		assert.NoError(t, err)

		secondWss, _, err := websocket.DefaultDialer.Dial(wsurl, header)
		assert.NoError(t, err)

		secondWss.SetCloseHandler(func(code int, text string) error {
			assert.Equal(t, code, websocket.ClosePolicyViolation)
			close(done)
			return err
		})

		go func() {
			for {
				_, _, err := secondWss.ReadMessage()
				if err != nil {
					break
				}
			}
		}()
		select {
		case <-time.After(timeout):
			assert.Fail(t, "Timeout. Expecting second connection to close")
			break
		case <-done:
			break
		}
	})

	t.Run("Should accept connections with same user id with different connection group", func(t *testing.T) {
		done := make(chan int)
		_, _, err := websocket.DefaultDialer.Dial(wsurl, http.Header{
			"X-User-ID":    []string{"1234"},
			"X-User-Group": []string{"viewer"},
		})

		assert.NoError(t, err)

		secondWss, _, err := websocket.DefaultDialer.Dial(wsurl, http.Header{
			"X-User-ID":    []string{"1234"},
			"X-User-Group": []string{"editor"},
		})

		assert.NoError(t, err)

		go func() {
			for {
				_, _, err := secondWss.ReadMessage()
				assert.NoError(t, err)
				if err != nil {
					close(done)
					break
				}
			}
		}()
		select {
		case <-time.After(timeout):
			assert.Fail(t, "Timeout. Expecting second connection to close")
			break
		case <-time.After(3 * time.Second):
			// Second connection is established and there is no error
			break
		case <-done:
			break
		}
	})

}
