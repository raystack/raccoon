package pubsub

import "fmt"

type unknownTopicError struct {
	Topic, Project string
}

func (e *unknownTopicError) Error() string {
	return fmt.Sprintf(
		`topic %q doesn't exist in %q project`, e.Topic, e.Project,
	)
}
