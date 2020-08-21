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
	format string
	m      *sync.Mutex
	topics map[string]string
}

func NewRouter() *Router {
	return &Router{
		format: config.NewTopicConfig().GetFormat(),
		m:      &sync.Mutex{},
		topics: make(map[string]string),
	}
}

func (r *Router) getTopic(eventType string) string {
	if !r.isExist(eventType) {
		r.m.Lock()
		r.topics[eventType] = fmt.Sprintf(r.format, eventType)
		r.m.Unlock()
	}
	return r.get(eventType)
}

func (r *Router) isExist(eventType string) bool {
	r.m.Lock()
	defer r.m.Unlock()
	return r.topics[eventType] != ""
}

func (r *Router) get(eventType string) string {
	r.m.Lock()
	defer r.m.Unlock()
	return r.topics[eventType]
}
