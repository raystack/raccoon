package kafka

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	pb "github.com/raystack/raccoon/proto"
	"github.com/raystack/raccoon/publisher"
)

func New() (*Kafka, error) {
	kp, err := newKafkaClient(config.Server.PublisherKafka.ToKafkaConfigMap())
	if err != nil {
		return &Kafka{}, err
	}
	k := &Kafka{
		kp:                  kp,
		flushInterval:       config.Server.PublisherKafka.FlushInterval,
		topicFormat:         config.Server.EventDistribution.PublisherPattern,
		deliveryChannelSize: config.Server.Worker.DeliveryChannelSize,
	}

	go k.ReportStats()

	return k, nil
}

func NewFromClient(client Client, flushInterval int, topicFormat string, deliveryChannelSize int) *Kafka {
	return &Kafka{
		kp:                  client,
		flushInterval:       flushInterval,
		topicFormat:         topicFormat,
		deliveryChannelSize: deliveryChannelSize,
	}
}

type Kafka struct {
	kp                  Client
	flushInterval       int
	topicFormat         string
	deliveryChannelSize int
}

// ProduceBulk messages to kafka. Block until all messages are sent. Return array of error. Order of Errors is guaranteed.
func (pr *Kafka) ProduceBulk(events []*pb.Event, connGroup string) error {
	errors := make([]error, len(events))
	totalProcessed := 0
	deliveryChannel := make(chan kafka.Event, pr.deliveryChannelSize)
	for order, event := range events {
		topic := strings.Replace(pr.topicFormat, "%s", event.Type, 1)
		message := &kafka.Message{
			Value:          event.EventBytes,
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Opaque:         order,
		}

		err := pr.kp.Produce(message, deliveryChannel)
		if err != nil {
			metrics.Increment(
				"kafka_messages_undelivered_total",
				map[string]string{
					"topic":      topic,
					"conn_group": connGroup,
					"event_type": event.Type,
				},
			)
			if err.Error() == "Local: Unknown topic" {
				errors[order] = fmt.Errorf("%v %s", err, topic)
				metrics.Increment(
					"kafka_unknown_topic_failure_total",
					map[string]string{
						"topic":      topic,
						"event_type": event.Type,
						"conn_group": connGroup,
					},
				)
			} else {
				errors[order] = err
			}
			continue
		}
		totalProcessed++
	}

	// Wait for deliveryChannel as many as processed
	for range totalProcessed {
		var (
			deliveryReport = <-deliveryChannel
			msg            = deliveryReport.(*kafka.Message)
			order          = msg.Opaque.(int)
			eventType      = events[order].Type
		)
		if msg.TopicPartition.Error != nil {
			metrics.Increment(
				"kafka_messages_undelivered_total",
				map[string]string{
					"topic":      *msg.TopicPartition.Topic,
					"conn_group": connGroup,
					"event_type": eventType,
				},
			)
			errors[order] = msg.TopicPartition.Error
			continue
		}
		metrics.Increment(
			"kafka_messages_delivered_total",
			map[string]string{
				"topic":      *msg.TopicPartition.Topic,
				"conn_group": connGroup,
				"event_type": eventType,
			},
		)
	}

	if cmp.Or(errors...) != nil {
		return &publisher.BulkError{Errors: errors}
	}

	return nil
}

func (pr *Kafka) ReportStats() {
	for v := range pr.kp.Events() {
		switch e := v.(type) {
		case *kafka.Stats:
			var stats map[string]interface{}
			json.Unmarshal([]byte(e.String()), &stats)

			brokers := stats["brokers"].(map[string]interface{})
			metrics.Gauge("kafka_tx_messages_total", stats["txmsgs"], map[string]string{})
			metrics.Gauge("kafka_tx_messages_bytes_total", stats["txmsg_bytes"], map[string]string{})
			for _, broker := range brokers {
				brokerStats := broker.(map[string]interface{})
				rttValue := brokerStats["rtt"].(map[string]interface{})
				nodeName := strings.Split(brokerStats["nodename"].(string), ":")[0]

				metrics.Gauge("kafka_brokers_tx_total", brokerStats["tx"], map[string]string{"broker": nodeName})
				metrics.Gauge("kafka_brokers_tx_bytes_total", brokerStats["txbytes"], map[string]string{"broker": nodeName})
				metrics.Gauge("kafka_brokers_rtt_average_milliseconds", rttValue["avg"], map[string]string{"broker": nodeName})
			}

		default:
			fmt.Printf("Ignored %v \n", e)
		}
	}
}

// Close wait for outstanding messages to be delivered within given flush interval timeout.
func (pr *Kafka) Close() error {
	remaining := pr.kp.Flush(pr.flushInterval)
	logger.Info(fmt.Sprintf("Outstanding events still un-flushed : %d", remaining))
	pr.kp.Close()
	if remaining > 0 {
		return &publisher.UnflushedEventsError{Count: remaining}
	}
	return nil
}

func (pr *Kafka) Name() string {
	return "kafka"
}

type ProducerStats struct {
	EventCounts map[string]int
	ErrorCounts map[string]int
}
