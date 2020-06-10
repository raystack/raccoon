package buffer

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer Produce data to kafka synchornously
type KafkaProducer interface {
	Produce(message *kafka.Message, deliveryChannel chan kafka.Event) error
}

// Worker spawn goroutine as much as PoolNumbers that will listen to In channel. On Close wait for all data in In channel to be processed.
type Worker struct {
	PoolNumbers   int
	In            <-chan []byte
	kafkaProducer KafkaProducer
	wg            sync.WaitGroup
}

// NewWorker create new Worker struct given poolNumber and max In channel buffer.
func NewWorker(poolNumber int, inChannel <-chan []byte, kafkaProducer KafkaProducer) Worker {
	return Worker{
		PoolNumbers:   poolNumber,
		In:            inChannel,
		kafkaProducer: kafkaProducer,
		wg:            sync.WaitGroup{},
	}
}

// StartWorker initialize worker pool as much as Worker.poolNumber
func (w *Worker) StartWorker() {
	w.wg.Add(w.PoolNumbers)
	for i := 0; i <= w.PoolNumbers; i++ {
		go func() {
			deliveryChan := make(chan kafka.Event, 1)
			for event := range w.In {
				message := kafka.Message{
					Value: event,
				}
				onFailRetry(&message, deliveryChan, w.kafkaProducer)
			}
			w.wg.Done()
		}()
	}
}

// Flush wait for remaining data to be processed. Call this after closing In channel
func (w *Worker) Flush() {
	w.wg.Wait()
}

func onFailRetry(message *kafka.Message, deliveryChan chan kafka.Event, producer KafkaProducer) {
	if err := producer.Produce(message, deliveryChan); err != nil {
		onFailRetry(message, deliveryChan, producer)
	}
}
