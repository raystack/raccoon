package metrics

import (
	"fmt"
	client "gopkg.in/alexcesaro/statsd.v2"
	"raccoon/config"
)

var instance *Statsd

type Statsd struct {
	c *client.Client
}

func (s *Statsd) Count(bucket string, i int, tags string) {
	s.c.Count(withTags(bucket, tags), i)
}

func (s *Statsd) Timing(bucket string, t int64, tags string) {
	s.c.Timing(withTags(bucket, tags), t)
}

func (s *Statsd) Increment(bucket string, tags string) {
	s.c.Increment(withTags(bucket, tags))
}

func (s *Statsd) Decrement(bucket string, tags string) {
	s.c.Count(withTags(bucket, tags), -1)
}

func (s *Statsd) Gauge(bucket string, val interface{}, tags string) {
	s.c.Gauge(withTags(bucket, tags), val)
}

func (s *Statsd) Close() {
	s.c.Close()
}

func withTags(bucket, tags string) string {
	return fmt.Sprintf("%v,%v", bucket, tags)
}

func Setup() {
	if instance == nil {
		c, err := client.New(
			client.Address(config.StatsdConfigLoader().Address),
			client.FlushPeriod(config.StatsdConfigLoader().FlushPeriod()))
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		instance = &Statsd{
			c: c,
		}

	}
}

func Count(bucket string, i int, tags string) {
	Setup()
	instance.Count(bucket, i, tags)
}

func Timing(bucket string, t int64, tags string) {
	Setup()
	instance.Timing(bucket, t, tags)
}

func Increment(bucket string, tags string) {
	Setup()
	instance.Increment(bucket, tags)
}

func Decrement(bucket string, tags string) {
	Setup()
	instance.Decrement(bucket, tags)
}

func Gauge(bucket string, val interface{}, tags string) {
	Setup()
	instance.Gauge(bucket, val, tags)
}

func Close() {
	Setup()
	instance.Close()
}

func Instance() *Statsd {
	Setup()
	return instance
}

func SetVoid() {
	c, _ := client.New(client.Mute(true))
	instance = &Statsd{
		c: c,
	}
}
