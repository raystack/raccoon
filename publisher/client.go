package publisher

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"raccoon/config"
	"raccoon/logger"
	"syscall"
)

func NewProducer(kp KafkaProducer, config config.KafkaConfig) *Producer {
	return &Producer{
		kp:     kp,
		Config: config,
	}
}

type Producer struct {
	kp               KafkaProducer
	Config           config.KafkaConfig
}

func (pr *Producer) Produce(msg *kafka.Message, deliveryChannel chan kafka.Event) error {

	produceErr := pr.kp.Produce(msg, deliveryChannel)

	if produceErr != nil {
		logger.Error("Producer failed to send message", produceErr)
		return produceErr
	}

	e := <-deliveryChannel
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logger.Error(fmt.Sprintf("Kafka message delivery failed.%s", m.TopicPartition.Error))
		return m.TopicPartition.Error
	} else {
		logger.Debug(fmt.Sprintf("Delivered message to topic %s [%d] at offset %s",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
	}
	return nil
}

func (pr *Producer) Close() {
	pr.kp.Close()
}

func (pr *Producer) Flush(flushInterval int) {
	 pr.kp.Flush(flushInterval)
}

func ShutdownProducer(ctx context.Context, pr *Producer) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Debug(fmt.Sprintf("[Kafka.Producer] Received a signal %s", sig))
			flushInterval := config.NewKafkaConfig().FlushInterval()
			logger.Debug(fmt.Sprintf("Wait %s ms for all messages to be delivered",flushInterval))
			pr.Flush(flushInterval)
			logger.Debug("Closing Producer")
			pr.Close()
		default:
			logger.Error(fmt.Sprintf("[Kafka.Producer] Received a unexpected signal %s", sig))
		}
	}
}
