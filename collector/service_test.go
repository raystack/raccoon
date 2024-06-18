package collector

import (
	"reflect"
	"testing"
)

func TestNewChannelCollector(t *testing.T) {
	type args struct {
		c chan CollectRequest
	}

	c := make(chan CollectRequest)
	tests := []struct {
		name string
		args args
		want Collector
	}{
		{
			name: "Get Collector",
			args: args{
				c: c,
			},
			want: &ChannelCollector{
				ch: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChannelCollector(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChannelCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}
