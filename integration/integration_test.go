package integration

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"math/rand"
	"net/http"
	"os"
	"source.golabs.io/mobile/clickstream-go-proto/gojek/clickstream/de"
	"testing"
	"time"
)

var uuid string
var timeout time.Duration
var topic string
var url string
var bootstrapServers string

func TestMain(m *testing.M) {
	uuid = fmt.Sprintf("%d-test", rand.Int())
	timeout = 120 * time.Second
	topic = os.Getenv("KAFKA_TOPIC")
	url = fmt.Sprintf("%v/api/v1/events", os.Getenv("INTEGTEST_HOST"))
	bootstrapServers = os.Getenv("INTEGTEST_HOST")
	m.Run()
}

func TestIntegration(t *testing.T) {
	accessToken, err := FetchAccessToken()
	assert.NoError(t, err)
	header := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %v", accessToken)},
		"GO-User-ID":    []string{"1234"},
	}
	t.Run("Should response with BadRequest when sending invalid request", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(url, header)

		if err != nil {
			assert.Fail(t, fmt.Sprintf("fail to connect. %v", err))
		}

		wss.WriteMessage(websocket.BinaryMessage, []byte{1})

		mType, resp, err := wss.ReadMessage()
		r := &de.EventResponse{}
		_ = proto.Unmarshal(resp, r)
		assert.Equal(t, mType, websocket.BinaryMessage)
		assert.Empty(t, err)
		assert.Equal(t, de.Status_ERROR, r.Status)
		assert.Equal(t, de.Code_BAD_REQUEST, r.Code)
		assert.NotEmpty(t, r.Reason)
		assert.Empty(t, r.Data)

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should response with success when request is processed", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(url, header)

		if err != nil {
			panic(err)
		}

		req := &de.EventRequest{
			ReqGuid:  "1234",
			SentTime: ptypes.TimestampNow(),
			Data: []*de.CSEventMessage{{
				EventGuid:      uuid,
				EventType:      "ride",
				EventName:      "ride",
				EventTimestamp: ptypes.TimestampNow()},
			},
		}
		bReq, _ := proto.Marshal(req)
		wss.WriteMessage(websocket.BinaryMessage, bReq)

		mType, resp, err := wss.ReadMessage()
		r := &de.EventResponse{}
		_ = proto.Unmarshal(resp, r)
		assert.Equal(t, mType, websocket.BinaryMessage)
		assert.Empty(t, err)
		assert.Equal(t, r.Code, de.Code_OK)
		assert.Equal(t, r.Status, de.Status_SUCCESS)
		assert.Equal(t, r.Reason, "")
		assert.Equal(t, r.Data, map[string]string{"req_guid": "1234"})

		wss.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(100*time.Millisecond))
	})

	t.Run("Should be able to consume published message", func(t *testing.T) {
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
			"group.id":          "my-local-group",
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			assert.Fail(t, "setup kafka consumer failed")
		}

		e := c.Subscribe(topic, nil)
		if e != nil {
			assert.Fail(t, fmt.Sprintf("Pls try again. %v", e))
		}

		for {
			select {
			case <-time.After(timeout):
				assert.Fail(t, "timeout")
				return
			default:
				msg, err := c.ReadMessage(timeout)
				if err != nil {
					continue
				}
				m := &de.CSEventMessage{}
				proto.Unmarshal(msg.Value, m)
				if m.EventGuid == uuid {
					return
				}
			}
		}

	})

	t.Run("Should close connection when client is unresponsive", func(t *testing.T) {
		wss, _, err := websocket.DefaultDialer.Dial(url, header)

		if err != nil {
			assert.Fail(t, err.Error())
		}

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
		_, _, err := websocket.DefaultDialer.Dial(url, header)

		assert.NoError(t, err)

		secondWss, _, err := websocket.DefaultDialer.Dial(url, header)

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

}
