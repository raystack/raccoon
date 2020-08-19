package worker

import (
	"context"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"raccoon/config"
	"sync"
)

type TopicsCreator interface {
	CreateTopics(context.Context, []kafka.TopicSpecification, ...kafka.CreateTopicsAdminOption) ([]kafka.TopicResult, error)
}

// Router create and cache topic based on provided format.
type Router struct {
	topicsCreator     TopicsCreator
	format            string
	topicConfigMap    map[string]string
	numPartitions     int
	replicationFactor int
	m                 *sync.Mutex
	topics            map[string]string
}

func NewRouter(creator TopicsCreator) *Router {
	return &Router{
		topicsCreator:     creator,
		format:            config.NewTopicConfig().GetFormat(),
		topicConfigMap:    config.NewTopicConfig().ToTopicConfigMap(),
		numPartitions:     config.NewTopicConfig().NumPartitions,
		replicationFactor: config.NewTopicConfig().GetReplicationFactor(),
		m:                 &sync.Mutex{},
		topics:            make(map[string]string),
	}
}

func (r *Router) getTopic(eventType string) (string, error) {
	if !r.isExist(eventType) {
		ctx := context.Background()
		topicResults, err := r.topicsCreator.CreateTopics(ctx, []kafka.TopicSpecification{{
			Topic:             fmt.Sprintf(r.format, eventType),
			NumPartitions:     r.numPartitions,
			ReplicationFactor: r.replicationFactor,
			Config:            r.topicConfigMap,
		}})

		for _, res := range topicResults {
			if res.Error.Code() != kafka.ErrNoError && res.Error.Code() != kafka.ErrTopicAlreadyExists {
				return "", res.Error
			}
			r.add(eventType)
		}
		if err != nil {
			return "", err
		}
	}
	return r.get(eventType), nil
}

func (r *Router) isExist(eventType string) bool {
	r.m.Lock()
	defer r.m.Unlock()
	return r.topics[eventType] != ""
}

func (r *Router) add(eventType string) {
	r.m.Lock()
	r.topics[eventType] = fmt.Sprintf(r.format, eventType)
	r.m.Unlock()
}

func (r *Router) get(eventType string) string {
	r.m.Lock()
	defer r.m.Unlock()
	return r.topics[eventType]
}
