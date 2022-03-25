package publisher

import (
	"encoding/json"
	"fmt"
	pb "raccoon/proto"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	// Importing librd to make it work on vendor mode
	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	"strings"

	_ "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka/librdkafka"
)

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	// ProduceBulk message to kafka. Block until all messages are sent. Return array of error. Order is not guaranteed.
	ProduceBulk(events []*pb.Event, deliveryChannel chan kafka.Event) error
}

func NewKafka() (*Kafka, error) {
	kp, err := newKafkaClient(config.PublisherKafka.ToKafkaConfigMap())
	if err != nil {
		return &Kafka{}, err
	}
	return &Kafka{
		kp:            kp,
		flushInterval: config.PublisherKafka.FlushInterval,
		topicFormat:   config.EventDistribution.PublisherPattern,
	}, nil
}

func NewKafkaFromClient(client Client, flushInterval int, topicFormat string) *Kafka {
	return &Kafka{
		kp:            client,
		flushInterval: flushInterval,
		topicFormat:   topicFormat,
	}
}

type Kafka struct {
	kp            Client
	flushInterval int
	topicFormat   string
}

// ProduceBulk messages to kafka. Block until all messages are sent. Return array of error. Order of Errors is guaranteed.
// DeliveryChannel needs to be exclusive. DeliveryChannel is exposed for recyclability purpose.
func (pr *Kafka) ProduceBulk(events []*pb.Event, deliveryChannel chan kafka.Event) error {
	errors := make([]*EventError, len(events))
	totalProcessed := 0
	for order, event := range events {
		topic := fmt.Sprintf(pr.topicFormat, event.Type)
		message := &kafka.Message{
			Value:          event.EventBytes,
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Opaque:         order,
		}

		err := pr.kp.Produce(message, deliveryChannel)
		if err != nil {
			if err.Error() == "Local: Unknown topic" {
				errors[order] = &EventError{Err: fmt.Errorf("%v %s", err, topic), EventType: event.Type}
				metrics.Increment("kafka_unknown_topic_failure_total", fmt.Sprintf("topic=%s,event_type=%s", topic, event.Type))
			} else {
				errors[order] = &EventError{Err: err, EventType: event.Type}
			}
			continue
		}
		totalProcessed++
	}
	// Wait for deliveryChannel as many as processed
	for i := 0; i < totalProcessed; i++ {
		d := <-deliveryChannel
		m := d.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			order := m.Opaque.(int)
			errors[order] = &EventError{Err: m.TopicPartition.Error, EventType: events[order].Type}
		}
	}

	if allNil(errors) {
		return nil
	}
	return BulkError{Errors: errors}
}

func (pr *Kafka) ReportStats() {
	for v := range pr.kp.Events() {
		switch e := v.(type) {
		case *kafka.Stats:
			var stats map[string]interface{}
			json.Unmarshal([]byte(e.String()), &stats)

			brokers := stats["brokers"].(map[string]interface{})
			metrics.Gauge("kafka_tx_messages_total", stats["txmsgs"], "")
			metrics.Gauge("kafka_tx_messages_bytes_total", stats["txmsg_bytes"], "")
			for _, broker := range brokers {
				brokerStats := broker.(map[string]interface{})
				rttValue := brokerStats["rtt"].(map[string]interface{})
				nodeName := strings.Split(brokerStats["nodename"].(string), ":")[0]

				metrics.Gauge("kafka_brokers_tx_total", brokerStats["tx"], fmt.Sprintf("host=%s,broker=true", nodeName))
				metrics.Gauge("kafka_brokers_tx_bytes_total", brokerStats["txbytes"], fmt.Sprintf("host=%s,broker=true", nodeName))
				metrics.Gauge("kafka_brokers_rtt_average_milliseconds", rttValue["avg"], fmt.Sprintf("host=%s,broker=true", nodeName))
			}

		default:
			fmt.Printf("Ignored %v \n", e)
		}
	}
}

// Close wait for outstanding messages to be delivered within given flush interval timeout.
func (pr *Kafka) Close() int {
	remaining := pr.kp.Flush(pr.flushInterval)
	logger.Info(fmt.Sprintf("Outstanding events still un-flushed : %d", remaining))
	pr.kp.Close()
	return remaining
}

func allNil(errors []*EventError) bool {
	for _, eErr := range errors {
		if eErr != nil && eErr.Err != nil {
			return false
		}
	}
	return true
}

type BulkError struct {
	Errors []*EventError
}

type EventError struct {
	Err       error
	EventType string
}

func (b BulkError) Error() string {
	err := "error when sending messages: "
	for i, eErr := range b.Errors {
		if i != 0 {
			err += fmt.Sprintf(", %v", eErr.Err)
			continue
		}
		err += eErr.Err.Error()
	}
	return err
}
