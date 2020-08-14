package publisher

import (
	"encoding/json"
	"fmt"

	"raccoon/config"
	"raccoon/logger"
	"raccoon/metrics"
	"strings"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	_ "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka/librdkafka"
)

type Event struct {
	Datum []byte
	Topic string
}

// KafkaProducer Produce data to kafka synchronously
type KafkaProducer interface {
	// ProduceBulk message to kafka. Block until all messages are sent. Return array of error. Order is not guaranteed.
	ProduceBulk(events []Event, deliveryChannel chan kafka.Event) error
}

func NewKafka(config config.KafkaConfig) (*Kafka, error) {
	kp, err := newKafkaClient(config.ToKafkaConfigMap())
	if err != nil {
		return &Kafka{}, err
	}
	return &Kafka{
		kp:     kp,
		Config: config,
	}, nil
}

func NewKafkaFromClient(client Client, config config.KafkaConfig) *Kafka {
	return &Kafka{
		kp:     client,
		Config: config,
	}
}

type Kafka struct {
	kp     Client
	Config config.KafkaConfig
	topics map[string]string
}

// ProduceBulk messages to kafka. Block until all messages are sent. Return array of error. Order of Errors is guaranteed.
// DeliveryChannel needs to be exclusive. DeliveryChannel is exposed for recyclability purpose.
func (pr *Kafka) ProduceBulk(events []Event, deliveryChannel chan kafka.Event) error {
	errors := make([]error, len(events))
	totalProcessed := 0
	for order, event := range events {
		message := &kafka.Message{
			Value:          event.Datum,
			TopicPartition: kafka.TopicPartition{Topic: &event.Topic, Partition: kafka.PartitionAny},
			Opaque:         order,
		}

		err := pr.kp.Produce(message, deliveryChannel)
		if err != nil {
			errors[order] = err
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
			errors[order] = m.TopicPartition.Error
		}
	}

	if allNil(errors) {
		return nil
	}
	return BulkError{Errors: errors}
}

func (pr *Kafka) ReportStats() {
	for  v := range pr.kp.Events() {
		switch e := v.(type) {
		case *kafka.Stats:
			var stats map[string]interface{}
			json.Unmarshal([]byte(e.String()), &stats)

			brokers := stats["brokers"].(map[string]interface{})
			metrics.Gauge("kafka.total.produced", stats["txmsgs"], "")
			metrics.Gauge("kafka.total_bytes.produced", stats["txmsg_bytes"],"")
			for _, broker := range brokers {
				brokerStats := broker.(map[string]interface{})
				rttValue := brokerStats["rtt"].(map[string]interface{})
				nodeName := strings.Split(brokerStats["nodename"].(string),":")[0]

				metrics.Gauge("kafka.request.sent", brokerStats["tx"], fmt.Sprintf("host=%s,broker=true",nodeName))
				metrics.Gauge("kafka.bytes.sent", brokerStats["txbytes"], fmt.Sprintf("host=%s,broker=true",nodeName))
				metrics.Gauge("kafka.round_trip_time.ms", rttValue["avg"], fmt.Sprintf("host=%s,broker=true",nodeName))
			}

		default:
			fmt.Printf("Ignored %v \n", e)
		}
	}
}

// Close wait for outstanding messages to be delivered within given flush interval timeout.
func (pr *Kafka) Close() int {
	remaining := pr.kp.Flush(pr.Config.GetFlushInterval())
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
