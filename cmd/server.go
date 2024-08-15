package cmd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/raystack/raccoon/app"
	"github.com/raystack/raccoon/config"
	"github.com/raystack/raccoon/logger"
	"github.com/raystack/raccoon/metrics"
	"github.com/raystack/raccoon/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func serverCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "server",
		Short: "Start raccoon server",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Load()
			if err != nil {
				return err
			}
			middleware.Load()
			metrics.Setup()
			defer metrics.Close()
			logger.SetLevel(config.Log.Level)
			return app.Run()
		},
	}

	bindServerFlags(command)
	return command
}

func bindServerFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.SortFlags = false

	bindFlag(
		fs,
		&config.Log.Level,
		"LOG_LEVEL",
		"Level available are [debug info warn error fatal panic]",
	)
	bindFlag(
		fs,
		&config.Event.Ack,
		"EVENT_ACK",
		"Whether to send acknowledgements to clients or not. 1 to enable, 0 to disable.",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.AppPort,
		"SERVER_WEBSOCKET_PORT",
		"Port for the service to listen",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.ServerMaxConn,
		"SERVER_WEBSOCKET_MAX_CONN",
		"Maximum connection that can be handled by the server instance",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.ServerMaxConn,
		"SERVER_WEBSOCKET_READ_BUFFER_SIZE",
		"Input buffer size in bytes",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.WriteBufferSize,
		"SERVER_WEBSOCKET_WRITE_BUFFER_SIZE",
		"Output buffer size in bytes",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.Conn.IDHeader,
		"SERVER_WEBSOCKET_CONN_ID_HEADER",
		"Unique identifier for the server to maintain the connection",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.Conn.GroupHeader,
		"SERVER_WEBSOCKET_CONN_GROUP_HEADER",
		"Additional identifier for the server to maintain the connection",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.Conn.GroupDefault,
		"SERVER_WEBSOCKET_CONN_GROUP_DEFAULT",
		"Default connection group name",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.PingIntervalMS,
		"SERVER_WEBSOCKET_PING_INTERVAL_MS",
		"Interval of each ping to client in milliseconds",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.PongWaitIntervalMS,
		"SERVER_WEBSOCKET_PONG_WAIT_INTERVAL_MS",
		"Wait time for client to send Pong message in milliseconds",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.WriteWaitIntervalMS,
		"SERVER_WEBSOCKET_WRITE_WAIT_INTERVAL_MS",
		"Timeout deadline set on the writes in milliseconds",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.PingerSize,
		"SERVER_WEBSOCKET_PINGER_SIZE",
		"Number of goroutine spawned to ping clients",
	)
	bindFlag(
		fs,
		&config.Server.Websocket.CheckOrigin,
		"SERVER_WEBSOCKET_CHECK_ORIGIN",
		"Toggle CORS check on WebSocket API",
	)
	bindFlag(
		fs,
		&config.Server.Batch.DedupEnabled,
		"SERVER_BATCH_DEDUP_IN_CONNECTION_ENABLED",
		"Whether to discard duplicate messages",
	)
	bindFlag(
		fs,
		&config.Server.CORS.Enabled,
		"SERVER_CORS_ENABLED",
		"Toggle CORS check on REST API",
	)
	bindFlag(
		fs,
		&config.Server.CORS.AllowedOrigin,
		"SERVER_CORS_ALLOWED_ORIGIN",
		"Allowed origins for CORS. Use '*' to allow all",
	)
	bindFlag(
		fs,
		&config.Server.CORS.AllowedMethods,
		"SERVER_CORS_ALLOWED_METHODS",
		"Allowed HTTP Methods for CORS",
	)
	bindFlag(
		fs,
		&config.Server.CORS.AllowedHeaders,
		"SERVER_CORS_ALLOWED_HEADERS",
		"Allowed HTTP Headers for CORS",
	)
	bindFlag(
		fs,
		&config.Server.CORS.MaxAge,
		"SERVER_CORS_PREFLIGHT_MAX_AGE_SECONDS",
		"Max Age of preflight responses",
	)
	bindFlag(
		fs,
		&config.Worker.Buffer.ChannelSize,
		"WORKER_BUFFER_CHANNEL_SIZE",
		"Size of the buffer queue",
	)
	bindFlag(
		fs,
		&config.Worker.Buffer.FlushTimeoutMS,
		"WORKER_BUFFER_FLUSH_TIMEOUT_MS",
		"Timeout for flushing leftover messages on shutdown",
	)
	bindFlag(
		fs,
		&config.Worker.PoolSize,
		"WORKER_POOL_SIZE",
		"No of workers that processes the events concurrently",
	)
	bindFlag(
		fs,
		&config.Worker.DeliveryChannelSize,
		"WORKER_KAFKA_DELIVERY_CHANNEL_SIZE",
		"Delivery Channel size for Kafka publisher",
	)
	bindFlag(
		fs,
		&config.Event.DistributionPublisherPattern,
		"EVENT_DISTRIBUTION_PUBLISHER_PATTERN",
		"Topic template used for routing events",
	)
	bindFlag(
		fs,
		&config.Publisher.Type,
		"PUBLISHER_TYPE",
		"Publisher to use for transmitting events",
	)

	// kafka client dynamic configuration doesn't have corresponding
	// fields in configuration structs. So we use a placeholder reference
	// to these values.
	var placeholder string
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS",
		"Address of kafka brokers",
	)
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_ACKS",
		"Number of replica acknowledgement before kafka sends ack back to service",
	)
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_RETRIES",
		"Number of retries in case of failure",
	)
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_RETRY_BACKOFF_MS",
		"Backoff time on retry.",
	)
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_STATISTICS_INTERVAL_MS",
		"Interval of statistics emitted by kafka",
	)
	bindFlag(
		fs,
		&placeholder,
		"PUBLISHER_KAFKA_CLIENT_QUEUE_BUFFERING_MAX_MESSAGES",
		"Maximum number of messages allowed on the producer queue",
	)
	bindFlag(
		fs,
		&config.Publisher.Kafka.FlushInterval,
		"PUBLISHER_KAFKA_FLUSH_INTERVAL_MS",
		"Timeout for sending leftover messages on kafka publisher shutdown",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.CredentialsFile,
		"PUBLISHER_PUBSUB_CREDENTIALS",
		"Path to file containing GCP cloud credentials",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.ProjectId,
		"PUBLISHER_PUBSUB_PROJECT_ID",
		"Destination Google Cloud Project ID",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.TopicAutoCreate,
		"PUBLISHER_PUBSUB_TOPIC_AUTOCREATE",
		"Whether to create topic if it doesn't exist in PubSub",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.TopicRetentionPeriodMS,
		"PUBLISHER_PUBSUB_TOPIC_RETENTION_MS",
		"Retention period of created topics in milliseconds",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.PublishDelayThresholdMS,
		"PUBLISHER_PUBSUB_PUBLISH_DELAY_THRESHOLD_MS",
		"Maximum time to wait for before publishing a batch of events",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.PublishCountThreshold,
		"PUBLISHER_PUBSUB_PUBLISH_COUNT_THRESHOLD",
		"Maximum number of events to accumulate before transmission",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.PublishByteThreshold,
		"PUBLISHER_PUBSUB_PUBLISH_BYTE_THRESHOLD",
		"Maximum buffer size (in bytes)",
	)
	bindFlag(
		fs,
		&config.Publisher.PubSub.PublishTimeoutMS,
		"PUBLISHER_PUBSUB_PUBLISH_TIMEOUT_MS",
		"How long to wait before aborting a publish operation",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.Region,
		"PUBLISHER_KINESIS_AWS_REGION",
		"AWS Region of the target kinesis stream",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.CredentialsFile,
		"PUBLISHER_KINESIS_CREDENTIALS",
		"Path to file containing AWS credentials",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.StreamAutoCreate,
		"PUBLISHER_KINESIS_STREAM_AUTOCREATE",
		"Whether to create a stream if it doesn't exist in Kinesis",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.StreamMode,
		"PUBLISHER_KINESIS_STREAM_MODE",
		"Mode of auto-created streams. Valid values: [ON-DEMAND PROVISIONED]",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.DefaultShards,
		"PUBLISHER_KINESIS_STREAM_SHARDS",
		"Number of shards in auto-created streams",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.StreamProbeIntervalMS,
		"PUBLISHER_KINESIS_STREAM_PROBE_INTERVAL_MS",
		"time delay between stream status checks",
	)
	bindFlag(
		fs,
		&config.Publisher.Kinesis.PublishTimeoutMS,
		"PUBLISHER_KINESIS_PUBLISH_TIMEOUT_MS",
		"how long to wait for before aborting a publish operation",
	)
	bindFlag(
		fs,
		&config.Metric.RuntimeStatsRecordIntervalMS,
		"METRIC_RUNTIME_STATS_RECORD_INTERVAL_MS",
		"Time interval between runtime metric collection",
	)
	bindFlag(
		fs,
		&config.Metric.StatsD.Enabled,
		"METRIC_STATSD_ENABLED",
		"Enable statsd metric exporter",
	)
	bindFlag(
		fs,
		&config.Metric.StatsD.Address,
		"METRIC_STATSD_ADDRESS",
		"Address to reports the service metrics",
	)
	bindFlag(
		fs,
		&config.Metric.StatsD.FlushPeriodMS,
		"METRIC_STATSD_FLUSH_PERIOD_MS",
		"Interval for the service to push metrics",
	)
	bindFlag(
		fs,
		&config.Metric.Prometheus.Enabled,
		"METRIC_PROMETHEUS_ENABLED",
		"Enable prometheus http server to expose service metrics",
	)
	bindFlag(
		fs,
		&config.Metric.Prometheus.Path,
		"METRIC_PROMETHEUS_PATH",
		"The path at which prometheus server should serve metrics",
	)
	bindFlag(
		fs,
		&config.Metric.Prometheus.Port,
		"METRIC_PROMETHEUS_PORT",
		"Port to expose prometheus metrics on",
	)
}

func bindFlag(flag *pflag.FlagSet, ref any, name, desc string) {
	flagName := strings.ReplaceAll(
		strings.ToLower(name), "_", ".",
	)

	el := reflect.ValueOf(ref).Elem()
	kind := el.Kind()
	typ := el.Type()

	switch {
	case typ.Name() == "Duration":
		v := ref.(*time.Duration)
		flag.Var(durationFlag{v}, flagName, desc)
	case typ.Name() == "AckType":
		v := ref.(*config.AckType)
		flag.Var(ackTypeFlag{v}, flagName, desc)
	case kind == reflect.String:
		v := ref.(*string)
		flag.StringVar(v, flagName, *v, desc)
	case kind == reflect.Int:
		v := ref.(*int)
		flag.IntVar(v, flagName, *v, desc)
	case kind == reflect.Int64:
		v := ref.(*int64)
		flag.Int64Var(v, flagName, *v, desc)
	case kind == reflect.Uint32:
		v := ref.(*uint32)
		flag.Uint32Var(v, flagName, *v, desc)
	case kind == reflect.Bool:
		v := ref.(*bool)
		flag.BoolVar(v, flagName, *v, desc)
	case kind == reflect.Slice && typ.Elem().String() == "string":
		v := ref.(*[]string)
		flag.StringSliceVar(v, flagName, *v, desc)
	default:
		msg := fmt.Sprintf("unsupport flag. kind = %s, type = %s", kind, typ)
		panic(msg)
	}

	viper.BindPFlag(name, flag.Lookup(flagName))
}

type durationFlag struct {
	value *time.Duration
}

func (df durationFlag) String() string {
	if df.value == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *df.value/time.Millisecond)
}

func (df durationFlag) Set(raw string) error {
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}
	*df.value = time.Millisecond * time.Duration(v)
	return nil
}

func (df durationFlag) Type() string {
	return "int"
}

type ackTypeFlag struct {
	value *config.AckType
}

func (af ackTypeFlag) String() string {
	if af.value == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *af.value)
}

func (af ackTypeFlag) Set(raw string) error {
	v, err := strconv.ParseInt(raw, 10, 0)
	if err != nil {
		return fmt.Errorf("error parsing bool: %w", err)
	}
	*af.value = config.AckType(v)
	return nil
}

func (af ackTypeFlag) Type() string {
	return "int"
}
