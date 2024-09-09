package collector

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/raystack/raccoon/pkg/clock"
	"github.com/stretchr/testify/assert"
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
				ch:    c,
				clock: clock.Default,
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

func TestCollect(t *testing.T) {
	t.Run("It should mutate TimePushed to the time the collect request is acknowledged", func(t *testing.T) {
		now := time.Now()
		clk := &clock.Mock{}
		clk.On("Now").Return(now).Once()
		defer clk.AssertExpectations(t)

		ch := make(chan CollectRequest)
		defer close(ch)

		collector := &ChannelCollector{
			ch:    ch,
			clock: clk,
		}

		consumer := func(requests chan CollectRequest) {
			for range requests {
			}
		}
		go consumer(collector.ch)

		req := &CollectRequest{}
		assert.Nil(
			t, collector.Collect(context.Background(), req),
		)
		assert.Equal(t, req.TimePushed, now)
	})
}
