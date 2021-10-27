package rest

import (
	"net/http"
	"raccoon/pkg/collection"
	"raccoon/pkg/deserialization"
	"raccoon/pkg/serialization"
	"reflect"
	"testing"
)

func TestHandler_GetRESTAPIHandler(t *testing.T) {

	collector := &collection.MockCollector{}
	type args struct {
		c collection.Collector
	}
	tests := []struct {
		name string
		h    *Handler
		args args
		want http.HandlerFunc
	}{
		{
			name: "Return a REST API Handler",
			h:    &Handler{},
			args: args{
				c: collector,
			},
			want: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			if got := h.GetRESTAPIHandler(tt.args.c); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Handler.GetRESTAPIHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_getDeserializerSerializer(t *testing.T) {
	type args struct {
		contentType string
	}
	tests := []struct {
		name  string
		h     *Handler
		args  args
		want  deserialization.Deserializer
		want1 serialization.Serializer
	}{
		{
			name: "Return Proto Deserializer/Serializer",
			h:    &Handler{},
			args: args{
				contentType: "application/proto",
			},
			want:  deserialization.ProtoDeserilizer(),
			want1: serialization.ProtoDeserilizer(),
		},
		{
			name: "Return JSON Deserializer/Serializer",
			h:    &Handler{},
			args: args{
				contentType: "application/json",
			},
			want:  deserialization.JSONDeserializer(),
			want1: serialization.JSONSerializer(),
		},
		{
			name: "Return default Deserializer/Serializer",
			h:    &Handler{},
			args: args{
				contentType: "",
			},
			want:  deserialization.JSONDeserializer(),
			want1: serialization.JSONSerializer(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			got, got1 := h.getDeserializerSerializer(tt.args.contentType)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Handler.getDeserializerSerializer() got = %v, want %v", got, tt.want)
			}
			if reflect.TypeOf(got1) != reflect.TypeOf(tt.want1) {
				t.Errorf("Handler.getDeserializerSerializer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
