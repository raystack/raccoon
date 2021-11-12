package collection

import (
	"context"
	"reflect"
	"testing"
)

func TestNewChannelCollector(t *testing.T) {
	type args struct {
		c chan *CollectRequest
	}
	c := make(chan *CollectRequest)
	tests := []struct {
		name string
		args args
		want Collector
	}{
		{
			name: "Creating collector",
			args: args{
				c: c,
			},
			want: CollectFunction(func(ctx context.Context, req *CollectRequest) error {
				return nil
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelCollector(tt.args.c); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("NewChannelCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}
