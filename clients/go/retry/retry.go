package retry

import (
	"time"

	"github.com/raystack/raccoon/clients/go/log"
)

var (
	// Default retry configuration
	DefaultRetryMax  uint          = 3
	DefaultRetryWait time.Duration = 1 * time.Second
)

// Do retries the function f, waits time between retries until a successful call is made.
// wait is time to wait for the next call
// maxAttempts is the maximum number if calls to the function.
func Do(wait time.Duration, maxAttempts uint, f func() error) error {
	var calls uint = 0
	for {
		err := f()
		if err == nil {
			return nil
		}
		calls++
		if calls >= maxAttempts && maxAttempts != 0 {
			return err
		}
		time.Sleep(wait)
		log.Errorf("[Retry %v] Retrying after %v due to the Error: %v\n", calls, wait, err)
	}
}
