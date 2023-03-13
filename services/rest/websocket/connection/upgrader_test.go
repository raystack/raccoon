package connection

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/goto/raccoon/logger"
	"github.com/goto/raccoon/metrics"
	"github.com/stretchr/testify/assert"
)

type void struct{}

func (v void) Write(_ []byte) (int, error) {
	return 0, nil
}

func TestMain(t *testing.M) {
	logger.SetOutput(void{})
	metrics.SetVoid()
	os.Exit(t.Run())
}

var config = UpgraderConfig{
	ReadBufferSize:    10240,
	WriteBufferSize:   10240,
	CheckOrigin:       false,
	MaxUser:           2,
	PongWaitInterval:  time.Duration(60 * time.Second),
	WriteWaitInterval: time.Duration(5 * time.Second),
	ConnIDHeader:      "X-User-ID",
	ConnGroupHeader:   "",
	ConnGroupDefault:  "--default--",
}

func TestConnectionLifecycle(t *testing.T) {
	t.Run("Should increment total connection when upgraded", func(t *testing.T) {
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID": []string{"user1"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.NoError(t, u.err)
				assert.Equal(t, 1, upgrader.Table.TotalConnection())
			},
			onIteration: 1,
		})
	})

	t.Run("Should decrement total connection when client close the conn", func(t *testing.T) {
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID": []string{"user1"},
		}, {
			"X-User-ID": []string{"user1"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				if u.iteration == 1 {
					assert.Equal(t, 1, upgrader.Table.TotalConnection())
				}
				u.conn.Close()
			},
		})
	})
}

func TestConnectionGroup(t *testing.T) {
	t.Run("Should accept connections with same userid and different group", func(t *testing.T) {
		config.ConnGroupHeader = "X-User-Group"
		defer func() { config.ConnGroupHeader = "" }()
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}, {
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"editor"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.Equal(t, 2, upgrader.Table.TotalConnection())
				assert.NoError(t, u.err)
			},
			onIteration: 2,
		})
	})

	t.Run("Should use default when ConnGroupHeader is not provided", func(t *testing.T) {
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID": []string{"user1"},
		}, {
			"X-User-ID": []string{"user1"},
		}, {
			"X-User-ID": []string{"user1"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.EqualError(t, u.err, "disconnecting connection [--default--] user1: already connected")
			},
			onIteration: 3,
		})
	})

	t.Run("Should reject connections with same userid and same group", func(t *testing.T) {
		config.ConnGroupHeader = "X-User-Group"
		defer func() { config.ConnGroupHeader = "" }()
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}, {
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.Equal(t, 1, upgrader.Table.TotalConnection())
				assert.EqualError(t, u.err, "disconnecting connection [viewer] user1: already connected")
			},
			onIteration: 2,
		})
	})

	t.Run("Should be able to reconnect when connection is closed", func(t *testing.T) {
		config.ConnGroupHeader = "X-User-Group"
		defer func() { config.ConnGroupHeader = "" }()
		upgrader := NewUpgrader(config)
		headers := []http.Header{{
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}, {
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}, {
			"X-User-ID":    []string{"user1"},
			"X-User-Group": []string{"viewer"},
		}}
		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.Equal(t, 1, upgrader.Table.TotalConnection())
				assert.NoError(t, u.err)
				u.conn.Close()
			},
		})
	})
}

func TestConnectionRejection(t *testing.T) {
	t.Run("Should close new connection when max is reached", func(t *testing.T) {
		upgrader := NewUpgrader(config)
		headers := make([]http.Header, 0)
		for _, i := range []string{"1", "2", "3"} {
			headers = append(headers, http.Header{
				"X-User-ID": []string{"user-" + i},
			})
		}

		upgradeConnectionTestHelper(t, upgrader, headers, assertUpgrade{
			callback: func(u upgradeRes) {
				assert.EqualError(t, u.err, "max connection reached")
			},
			onIteration: 3,
		})
	})
}

// Prepare a websocket server with given upgrader and establish the connections with the given headers as many as given headers.
func upgradeConnectionTestHelper(t *testing.T, upgrader *Upgrader, headers []http.Header, f assertUpgrade) {
	res := make(chan upgradeRes)
	m := sync.Mutex{}
	iteration := 0
	r := mux.NewRouter()
	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		m.Lock()
		iteration++
		i := iteration
		m.Unlock()
		c, err := upgrader.Upgrade(rw, r)
		if f.onIteration == 0 {
			res <- upgradeRes{
				err:       err,
				iteration: i,
				conn:      c}
		}
		if i == f.onIteration {
			res <- upgradeRes{
				err:       err,
				iteration: i,
				conn:      c}
		}
	})
	server := httptest.NewServer(r)
	defer server.Close()
	connect := func(h http.Header) {
		websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), h)
	}
	for _, header := range headers {
		connect(header)
	}
	timeout := 5 * time.Second
	select {
	case <-time.After(timeout):
		t.Fatal("timeout, no error return from upgrader")
	case e := <-res:
		f.callback(e)
	}
}

// Struct to prepare upgrade assertion.
// If onIteration is provided, the assertion only run on the specified iteration of the passed headers. If onIteration is not provided or 0, assertion is run every upgrade.
type assertUpgrade struct {
	callback    func(u upgradeRes)
	onIteration int
}

type upgradeRes struct {
	err       error
	conn      Conn
	iteration int
}
