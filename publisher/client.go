package publisher

import (
	"clickstream-service/config"
	"clickstream-service/logger"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gopkg.in/errgo.v2/errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// check producer is already created or not, if already

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(ctx context.Context,kc config.KafkaConfig) (*Producer, error) {

	kafkaConfigMap := &kafka.ConfigMap{
		"bootstrap.servers": kc.BrokerList(),
	}

	producer, err := kafka.NewProducer(kafkaConfigMap)

	if err != nil {
		logger.Error(fmt.Sprintf("Error while creating new kafka producer. %s",err))
		return nil, errors.New(fmt.Sprintf("Error while creating new kafka producer. %s",err))
	}
	logger.Info("kafka producer created", producer)

	newProducer :=  &Producer{
		producer: producer,
	}
    go closeProducer(ctx, newProducer)
	return newProducer, nil
}

func closeProducer(ctx context.Context, p *Producer) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[Kafka.Producer] Received a signal %s", sig))
			logger.Info(fmt.Sprintf("[Kafka.Producer] waiting for 3 secs grace period before shutdown ", ))
			time.Sleep(3 * time.Second)
			logger.Info("Closing Producer")
			p.producer.Close()
		default:
			logger.Info(fmt.Sprintf("[Kafka.Producer] Received a unexpected signal %s", sig))
		}
	}
}

func (p *Producer) PublishMessage(msg, key []byte) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:      &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
		Key:   key,
	}, nil)

	if err != nil {
		logger.Error("Failed to publish message to kafka", err)
		return err
	}

	logger.Info("Message published to topic", p.topic)
	return nil
}
