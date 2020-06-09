package publisher

import (
	"clickstream-service/config"
	"clickstream-service/logger"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewProducer(kp KafkaProducer, config config.KafkaConfig) *Producer {
	return &Producer{
		kp:     kp,
		Config: config,
	}
}

type Producer struct {
	kp               KafkaProducer
	InflightMessages chan *kafka.Message
	Config           config.KafkaConfig
}

func (pr *Producer) Produce(msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)

	produceErr := pr.kp.Produce(msg, deliveryChan)

	if produceErr != nil {
		logger.Error("Kafka producer creation failed", produceErr)
		return produceErr
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.Error(fmt.Sprintf("Kafka message delivery failed.%s", m.TopicPartition.Error))
		return m.TopicPartition.Error
	} else {
		logger.Debug(fmt.Sprintf("Delivered message to topic %s [%d] at offset %s",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
	}
	close(deliveryChan)
	return nil
}

func shutdownProducer(ctx context.Context, pr *Producer) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[Kafka.Producer] Received a signal %s", sig))
			logger.Info(fmt.Sprintf("[Kafka.Producer] waiting for 3 secs grace period before shutdown "))
			time.Sleep(3 * time.Second)
			logger.Info("Closing Producer")
			pr.Close()
		default:
			logger.Info(fmt.Sprintf("[Kafka.Producer] Received a unexpected signal %s", sig))
		}
	}
}

func (pr *Producer) Close() {
	pr.kp.Close()
}
