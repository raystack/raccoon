package worker

import (
	"context"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"sync"
)

type TopicsCreator interface {
	CreateTopics(context.Context, []kafka.TopicSpecification, ...kafka.CreateTopicsAdminOption) ([]kafka.TopicResult, error)
}

// Router create and cache topic based on provided format.
type Router struct {
	topicsCreator TopicsCreator
	format        string
	m             *sync.Mutex
	topics        map[string]string
}

func NewRouter(creator TopicsCreator) *Router {
	// TODO: extract formatting to config
	return &Router{topicsCreator: creator, format: "clickstream-%s-log", topics: make(map[string]string), m: &sync.Mutex{}}
}

func (r *Router) getTopic(eventType string) (string, error) {
	if !r.isExist(eventType) {
		ctx := context.Background()
		topicResults, err := r.topicsCreator.CreateTopics(ctx, []kafka.TopicSpecification{{
			Topic: fmt.Sprintf(r.format, eventType),
			// TODO: extract these as configuration
			NumPartitions:     50,
			ReplicationFactor: 3,
			Config:            map[string]string{"retention.ms": "172800000"},
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
