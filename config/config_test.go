package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/raystack/raccoon/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	var testCases = []struct {
		Desc string
		Cfg  string
		Err  error
	}{
		{
			Desc: "should return an error if websocket.conn.id.header is not specified",
			Err: config.ConfigError{
				Env:  "SERVER_WEBSOCKET_CONN_ID_HEADER",
				Flag: "server.websocket.conn.id.header",
			},
		},
		{
			Desc: "should return an error if publisher type is pubsub and ProjectID is not specified",
			Cfg: heredoc.Doc(`
				server:
				  websocket:
				    conn:
				      id_header: "X-User-ID"
				publisher:
				  type: "pubsub"
			`),
			Err: config.ConfigError{
				Env:  "PUBLISHER_PUBSUB_PROJECT_ID",
				Flag: "publisher.pubsub.project.id",
			},
		},
		{
			Desc: "should return an error if publisher type is pubsub and credentials are not specified",
			Cfg: heredoc.Doc(`
				server:
				  websocket:
				    conn:
				      id_header: "X-User-ID"
				publisher:
				  type: "pubsub"
				  pubsub:
				    project_id: simulated-project-001
			`),
			Err: config.ConfigError{
				Env:  "PUBLISHER_PUBSUB_CREDENTIALS",
				Flag: "publisher.pubsub.credentials",
			},
		},
		{
			Desc: "should return an error if publisher type is kinesis and credentials are not specified",
			Cfg: heredoc.Doc(`
				server:
				  websocket:
				    conn:
				      id_header: "X-User-ID"
				publisher:
				  type: "kinesis"
			`),
			Err: config.ConfigError{
				Env:  "PUBLISHER_KINESIS_CREDENTIALS",
				Flag: "publisher.kinesis.credentials",
			},
		},
		{
			Desc: "should return an error if publisher type is kafka and bootstrap servers are not specified",
			Cfg: heredoc.Doc(`
				server:
				  websocket:
				    conn:
				      id_header: "X-User-ID"
				publisher:
				  type: "kafka"
			`),
			Err: config.ConfigError{
				Env:  "PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS",
				Flag: "publisher.kafka.client.bootstrap.servers",
			},
		},
		{
			Desc: "should return an error if an unknown publisher type is specified",
			Cfg: heredoc.Doc(`
				server:
				  websocket:
				    conn:
				      id_header: "X-User-ID"
				publisher:
				  type: "non-existent"
			`),
			Err: fmt.Errorf("unknown publisher: non-existent"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			fd, err := newTempFile()
			if err != nil {
				t.Errorf("error creating temporary file: %v", err)
				return
			}
			defer fd.Close()

			_, err = fmt.Fprint(fd, testCase.Cfg)
			if err != nil {
				t.Errorf("error writing test config: %v", err)
				return
			}

			err = config.Load(fd.Name())
			assert.Equal(t, err, testCase.Err)
		})
	}
}

type tempFile struct {
	*os.File
}

func (f tempFile) Close() error {
	f.File.Close()
	return os.Remove(f.File.Name())
}

func newTempFile() (tempFile, error) {
	fd, err := os.CreateTemp("", "raccoon-test-*")
	if err != nil {
		return tempFile{}, err
	}
	return tempFile{fd}, nil
}
