package metrics

import (
	"fmt"

	"github.com/odpf/raccoon/config"
	"github.com/odpf/raccoon/logger"
	client "gopkg.in/alexcesaro/statsd.v2"
)

var instance *Statsd

type Statsd struct {
	c *client.Client
}

func (s *Statsd) count(bucket string, i int, tags string) {
	s.c.Count(withTags(bucket, tags), i)
}

func (s *Statsd) timing(bucket string, t int64, tags string) {
	s.c.Timing(withTags(bucket, tags), t)
}

func (s *Statsd) increment(bucket string, tags string) {
	s.c.Increment(withTags(bucket, tags))
}

func (s *Statsd) decrement(bucket string, tags string) {
	s.c.Count(withTags(bucket, tags), -1)
}

func (s *Statsd) gauge(bucket string, val interface{}, tags string) {
	s.c.Gauge(withTags(bucket, tags), val)
}

func (s *Statsd) Close() {
	s.c.Close()
}

func withTags(bucket, tags string) string {
	return fmt.Sprintf("%v,%v", bucket, tags)
}

func Setup() error {
	if instance == nil {
		c, err := client.New(
			client.Address(config.MetricStatsd.Address),
			client.FlushPeriod(config.MetricStatsd.FlushPeriodMs))
		if err != nil {
			logger.Errorf("StatsD Set up failed to create client: %s", err.Error())
			return err
		}

		instance = &Statsd{
			c: c,
		}
	}
	return nil
}

func Count(bucket string, i int, tags string) {
	err := Setup()
	if err == nil {
		instance.count(bucket, i, tags)
	}
}

func Timing(bucket string, t int64, tags string) {
	err := Setup()
	if err == nil {
		instance.timing(bucket, t, tags)
	}
}

func Increment(bucket string, tags string) {
	err := Setup()
	if err == nil {
		instance.increment(bucket, tags)
	}
}

func Decrement(bucket string, tags string) {
	err := Setup()
	if err == nil {
		instance.decrement(bucket, tags)
	}
}

func Gauge(bucket string, val interface{}, tags string) {
	err := Setup()
	if err == nil {
		instance.gauge(bucket, val, tags)
	}
}

func Close() {
	err := Setup()
	if err == nil {
		instance.Close()
	}
}

func SetVoid() {
	c, _ := client.New(client.Mute(true))
	instance = &Statsd{
		c: c,
	}
}
