package worker

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestRouter(t *testing.T) {
	t.Run("Should return topic according to format", func(t *testing.T) {
		router := Router{
			m:      &sync.Mutex{},
			format: "prefix_%s_suffix",
			topics: make(map[string]string),
		}

		topic := router.getTopic("topic")
		assert.Equal(t, "prefix_topic_suffix", topic)
	})

	t.Run("isExist should return true when topic exist", func(t *testing.T) {
		router := Router{
			m:      &sync.Mutex{},
			format: "p_%s_s",
			topics: make(map[string]string),
		}
		assert.False(t, router.isExist("topic"))
		assert.Equal(t, "p_topic_s", router.getTopic("topic"))
		assert.True(t, router.isExist("topic"))
	})
}
