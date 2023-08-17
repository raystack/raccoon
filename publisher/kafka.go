package publisher

import (
	"encoding/json"
	"fmt"
	"strings"

	pb "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/raccoon/v1beta1"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/goto/raccoon/config"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
)

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	// ProduceBulk message to kafka. Block until all messages are sent. Return array of error. Order is not guaranteed.
	ProduceBulk(events []*pb.Event, connGroup string, deliveryChannel chan kafka.Event) error
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
func (pr *Kafka) ProduceBulk(events []*pb.Event, connGroup string, deliveryChannel chan kafka.Event) error {
	errors := make([]error, len(events))
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
			metrics.Increment("kafka_messages_delivered_total", fmt.Sprintf("success=false,conn_group=%s,event_type=%s", connGroup, event.Type))
			if err.Error() == "Local: Unknown topic" {
				errors[order] = fmt.Errorf("%v %s", err, topic)
				metrics.Increment("kafka_unknown_topic_failure_total", fmt.Sprintf("topic=%s,event_type=%s,conn_group=%s", topic, event.Type, connGroup))
			} else {
				errors[order] = err
			}
			continue
		}
		metrics.Increment("kafka_messages_delivered_total", fmt.Sprintf("success=true,conn_group=%s,event_type=%s", connGroup, event.Type))
		totalProcessed++
	}
	// Wait for deliveryChannel as many as processed
	for i := 0; i < totalProcessed; i++ {
		d := <-deliveryChannel
		m := d.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			eventType := events[i].Type
			metrics.Decrement("kafka_messages_delivered_total", fmt.Sprintf("success=true,conn_group=%s,event_type=%s", connGroup, eventType))
			metrics.Increment("kafka_messages_delivered_total", fmt.Sprintf("success=false,conn_group=%s,event_type=%s", connGroup, eventType))
			order := m.Opaque.(int)
			errors[order] = m.TopicPartition.Error
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

				metrics.Gauge("kafka_brokers_tx_total", brokerStats["tx"], fmt.Sprintf("broker=%s", nodeName))
				metrics.Gauge("kafka_brokers_tx_bytes_total", brokerStats["txbytes"], fmt.Sprintf("broker=%s", nodeName))
				metrics.Gauge("kafka_brokers_rtt_average_milliseconds", rttValue["avg"], fmt.Sprintf("broker=%s", nodeName))
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

func allNil(errors []error) bool {
	for _, err := range errors {
		if err != nil {
			return false
		}
	}
	return true
}

type ProducerStats struct {
	EventCounts map[string]int
	ErrorCounts map[string]int
}

type BulkError struct {
	Errors []error
}

func (b BulkError) Error() string {
	err := "error when sending messages: "
	for i, mErr := range b.Errors {
		if i != 0 {
			err += fmt.Sprintf(", %v", mErr)
			continue
		}
		err += mErr.Error()
	}
	return err
}
