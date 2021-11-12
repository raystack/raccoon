package rest

import (
	"net/http"
	"raccoon/collection"
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
