package websocket

import (
	"raccoon/http/websocket/connection"
	"raccoon/logger"
	"raccoon/metrics"
	"raccoon/pkg/deserialization"
	"raccoon/pkg/serialization"
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}
func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	metrics.SetVoid()
}

func TestNewHandler(t *testing.T) {
	type args struct {
		pingC chan connection.Conn
	}

	ugConfig := connection.UpgraderConfig{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       false,
		MaxUser:           100,
		PongWaitInterval:  60,
		WriteWaitInterval: 60,
		ConnIDHeader:      "x-conn-id",
		ConnGroupHeader:   "x-group",
	}
	pingC := make(chan connection.Conn)
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "creating a new handler",
			args: args{
				pingC: pingC,
			},
			want: &Handler{
				upgrader:    connection.NewUpgrader(ugConfig),
				PingChannel: pingC,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.pingC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Table(t *testing.T) {
	table := &connection.Table{}
	type fields struct {
		upgrader    *connection.Upgrader
		PingChannel chan connection.Conn
	}
	tests := []struct {
		name   string
		fields fields
		want   *connection.Table
	}{
		{
			name: "return table",
			fields: fields{
				upgrader: &connection.Upgrader{
					Table: table,
				},
			},
			want: table,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				upgrader:    tt.fields.upgrader,
				PingChannel: tt.fields.PingChannel,
			}
			if got := h.Table(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_getDeserializerSerializer(t *testing.T) {
	type fields struct {
		upgrader    *connection.Upgrader
		PingChannel chan connection.Conn
	}
	type args struct {
		messageType int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   deserialization.Deserializer
		want1  serialization.Serializer
	}{
		{
			name:   "get ProtoDeserializer",
			fields: fields{},
			args: args{
				messageType: websocket.BinaryMessage,
			},
			want:  deserialization.ProtoDeserilizer(),
			want1: serialization.ProtoDeserilizer(),
		},
		{
			name:   "get ProtoDeserializer",
			fields: fields{},
			args: args{
				messageType: websocket.BinaryMessage,
			},
			want:  deserialization.ProtoDeserilizer(),
			want1: serialization.ProtoDeserilizer(),
		},
		{
			name:   "get JSONDeserializer",
			fields: fields{},
			args: args{
				messageType: websocket.TextMessage,
			},
			want:  deserialization.JSONDeserializer(),
			want1: serialization.JSONSerializer(),
		},
		{
			name:   "get Default Deserializer",
			fields: fields{},
			args: args{
				messageType: websocket.TextMessage,
			},
			want:  deserialization.ProtoDeserilizer(),
			want1: serialization.ProtoDeserilizer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				upgrader:    tt.fields.upgrader,
				PingChannel: tt.fields.PingChannel,
			}
			got, got1 := h.getDeserializerSerializer(tt.args.messageType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.getDeserializerSerializer() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Handler.getDeserializerSerializer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
