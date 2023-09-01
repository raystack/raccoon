package metrics

import (
	"fmt"
	"strings"

	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	client "gopkg.in/alexcesaro/statsd.v2"
)

type Statsd struct {
	c *client.Client
}

func initStatsd() (*Statsd, error) {

	c, err := client.New(
		client.Address(config.MetricStatsd.Address),
		client.FlushPeriod(config.MetricStatsd.FlushPeriodMs))
	if err != nil {
		logger.Errorf("StatsD Set up failed to create client: %s", err.Error())
		return nil, err
	}
	return &Statsd{
		c: c,
	}, nil
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

func (s *Statsd) Count(metricName string, count int64, labels map[string]string) error {
	err := Setup()
	if err != nil {
		return err
	}
	tags := translateLabelIntoTags(labels)
	s.count(metricName, int(count), tags)
	return nil

}

func (s *Statsd) Histogram(metricName string, value int64, labels map[string]string) error {
	err := Setup()
	if err != nil {
		return err
	}
	tags := translateLabelIntoTags(labels)
	s.timing(metricName, value, tags)
	return nil
}

func (s *Statsd) Increment(bucket string, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.increment(bucket, tags)
	return nil
}

// func Decrement(bucket string, tags string) {
// 	err := Setup()
// 	if err == nil {
// 		instance.decrement(bucket, tags)
// 	}
// }

func (s *Statsd) Gauge(metricName string, value int64, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.gauge(metricName, value, tags)
	return nil
}

func setStatsDVoid() {
	c, _ := client.New(client.Mute(true))
	Instrument = &Statsd{
		c: c,
	}
}

func translateLabelIntoTags(labels map[string]string) string {
	labelArr := make([]string, len(labels))
	for key, value := range labels {
		labelArr = append(labelArr, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(labelArr, ",")
}
