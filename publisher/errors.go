package publisher

import "fmt"

type BulkError struct {
	Errors []error
}

func (b *BulkError) Error() string {
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

type UnflushedEventsError struct {
	Count int
}

func (e *UnflushedEventsError) Error() string {
	return fmt.Sprintf("%d events were not flushed", e.Count)
}
