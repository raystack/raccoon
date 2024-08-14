package metrics

import (
	"fmt"
	"strings"
	"time"

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
		client.FlushPeriod(time.Duration(config.MetricStatsd.FlushPeriodMS)*time.Millisecond))
	if err != nil {
		logger.Errorf("StatsD Set up failed to create client: %s", err.Error())
		return nil, err
	}
	return &Statsd{
		c: c,
	}, nil
}

func (s *Statsd) Close() {
	s.c.Close()
}

func withTags(bucket, tags string) string {
	return fmt.Sprintf("%v,%v", bucket, tags)
}

func (s *Statsd) Count(metricName string, count int64, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.c.Count(withTags(metricName, tags), int(count))
	return nil
}

func (s *Statsd) Histogram(metricName string, value int64, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.c.Timing(withTags(metricName, tags), value)
	return nil
}

func (s *Statsd) Increment(metricName string, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.c.Increment(withTags(metricName, tags))
	return nil
}

func (s *Statsd) Gauge(metricName string, value interface{}, labels map[string]string) error {
	tags := translateLabelIntoTags(labels)
	s.c.Gauge(withTags(metricName, tags), value)
	return nil
}

func translateLabelIntoTags(labels map[string]string) string {
	labelArr := make([]string, len(labels))
	for key, value := range labels {
		labelArr = append(labelArr, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(labelArr, ",")
}
