package config

import (
	"reflect"
	"strings"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

var Publisher publisher

type publisherPubSub struct {
	ProjectId               string `mapstructure:"project_id" cmdx:"publisher.pubsub.project.id"`
	TopicAutoCreate         bool   `mapstructure:"topic_autocreate" cmdx:"publisher.pubsub.topic.autocreate" default:"false"`
	TopicRetentionPeriodMS  int64  `mapstructure:"topic_retention_ms" cmdx:"publisher.pubsub.topic.retention.ms" default:"0"`
	PublishTimeoutMS        int64  `mapstructure:"publish_timeout_ms" cmdx:"publisher.pubsub.publish.timeout.ms" default:"60000"`
	PublishDelayThresholdMS int64  `mapstructure:"publish_delay_threshold_ms" cmdx:"publisher.pubsub.publish.delay.threshold.ms" default:"10"`
	PublishCountThreshold   int    `mapstructure:"publish_count_threshold" cmdx:"publisher.pubsub.publish.count.threshold" default:"100"`
	PublishByteThreshold    int    `mapstructure:"publish_byte_threshold" cmdx:"publisher.pubsub.publish.byte.threshold" default:"1000000"`
	CredentialsFile         string `mapstructure:"credentials" cmdx:"publisher.pubsub.credentials"`
}

type publisherKinesis struct {
	Region                string `mapstructure:"aws_region" cmdx:"publisher.kinesis.aws.region"`
	CredentialsFile       string `mapstructure:"credentials" cmdx:"publisher.kinesis.credentials"`
	StreamAutoCreate      bool   `mapstructure:"stream_autocreate" cmdx:"publisher.kinesis.stream.autocreate" default:"false"`
	StreamProbeIntervalMS int64  `mapstructure:"stream_probe_interval_ms" cmdx:"publisher.kinesis.stream.probe.interval.ms" default:"1000"`
	StreamMode            string `mapstructure:"stream_mode" cmdx:"publisher.kinesis.stream.mode" default:"ON_DEMAND"`
	DefaultShards         uint32 `mapstructure:"stream_shards" cmdx:"publisher.kinesis.stream.shards" default:"4"`
	PublishTimeoutMS      int64  `mapstructure:"publish_timeout_ms" cmdx:"publisher.kinesis.publish.timeout.ms" default:"60000"`
}

type kafkaClientConfig struct {
	BuiltinFeatures                     string `mapstructure:"builtin_features" cmdx:"publisher.kafka.client.builtin.features" default:"gzip,snappy,ssl,sasl,regex,lz4,sasl_plain,sasl_scram,plugins,zstd,sasl_oauthbearer,http,oidc"`
	ClientID                            string `mapstructure:"client_id" cmdx:"publisher.kafka.client.client.id" default:"rdkafka"`
	MetadataBrokerList                  string `mapstructure:"metadata_broker_list" cmdx:"publisher.kafka.client.metadata.broker.list" default:""`
	BootstrapServers                    string `mapstructure:"bootstrap_servers" cmdx:"publisher.kafka.client.bootstrap.servers" default:""`
	MessageMaxBytes                     string `mapstructure:"message_max_bytes" cmdx:"publisher.kafka.client.message.max.bytes" default:"1000000"`
	MessageCopyMaxBytes                 string `mapstructure:"message_copy_max_bytes" cmdx:"publisher.kafka.client.message.copy.max.bytes" default:"65535"`
	ReceiveMessageMaxBytes              string `mapstructure:"receive_message_max_bytes" cmdx:"publisher.kafka.client.receive.message.max.bytes" default:"100000000"`
	MaxInFlightRequestsPerConnection    string `mapstructure:"max_in_flight_requests_per_connection" cmdx:"publisher.kafka.client.max.in.flight.requests.per.connection" default:"1000000"`
	MaxInFlight                         string `mapstructure:"max_in_flight" cmdx:"publisher.kafka.client.max.in.flight" default:"1000000"`
	TopicMetadataRefreshIntervalMS      string `mapstructure:"topic_metadata_refresh_interval_ms" cmdx:"publisher.kafka.client.topic.metadata.refresh.interval.ms" default:"300000"`
	MetadataMaxAgeMS                    string `mapstructure:"metadata_max_age_ms" cmdx:"publisher.kafka.client.metadata.max.age.ms" default:"900000"`
	TopicMetadataRefreshFastIntervalMS  string `mapstructure:"topic_metadata_refresh_fast_interval_ms" cmdx:"publisher.kafka.client.topic.metadata.refresh.fast.interval.ms" default:"250"`
	TopicMetadataRefreshSparse          string `mapstructure:"topic_metadata_refresh_sparse" cmdx:"publisher.kafka.client.topic.metadata.refresh.sparse" default:"true"`
	TopicMetadataPropagationMaxMS       string `mapstructure:"topic_metadata_propagation_max_ms" cmdx:"publisher.kafka.client.topic.metadata.propagation.max.ms" default:"30000"`
	TopicBlacklist                      string `mapstructure:"topic_blacklist" cmdx:"publisher.kafka.client.topic.blacklist" default:""`
	SocketTimeoutMS                     string `mapstructure:"socket_timeout_ms" cmdx:"publisher.kafka.client.socket.timeout.ms" default:"60000"`
	SocketSendBufferBytes               string `mapstructure:"socket_send_buffer_bytes" cmdx:"publisher.kafka.client.socket.send.buffer.bytes" default:"0"`
	SocketReceiveBufferBytes            string `mapstructure:"socket_receive_buffer_bytes" cmdx:"publisher.kafka.client.socket.receive.buffer.bytes" default:"0"`
	SocketKeepaliveEnable               string `mapstructure:"socket_keepalive_enable" cmdx:"publisher.kafka.client.socket.keepalive.enable" default:"false"`
	SocketNagleDisable                  string `mapstructure:"socket_nagle_disable" cmdx:"publisher.kafka.client.socket.nagle.disable" default:"false"`
	SocketMaxFails                      string `mapstructure:"socket_max_fails" cmdx:"publisher.kafka.client.socket.max.fails" default:"1"`
	BrokerAddressTtl                    string `mapstructure:"broker_address_ttl" cmdx:"publisher.kafka.client.broker.address.ttl" default:"1000"`
	BrokerAddressFamily                 string `mapstructure:"broker_address_family" cmdx:"publisher.kafka.client.broker.address.family" default:"any"`
	SocketConnectionSetupTimeoutMS      string `mapstructure:"socket_connection_setup_timeout_ms" cmdx:"publisher.kafka.client.socket.connection.setup.timeout.ms" default:"30000"`
	ConnectionsMaxIdleMS                string `mapstructure:"connections_max_idle_ms" cmdx:"publisher.kafka.client.connections.max.idle.ms" default:"0"`
	ReconnectBackoffMS                  string `mapstructure:"reconnect_backoff_ms" cmdx:"publisher.kafka.client.reconnect.backoff.ms" default:"100"`
	ReconnectBackoffMaxMS               string `mapstructure:"reconnect_backoff_max_ms" cmdx:"publisher.kafka.client.reconnect.backoff.max.ms" default:"10000"`
	StatisticsIntervalMS                string `mapstructure:"statistics_interval_ms" cmdx:"publisher.kafka.client.statistics.interval.ms" default:"0"`
	LogQueue                            string `mapstructure:"log_queue" cmdx:"publisher.kafka.client.log.queue" default:"false"`
	LogThreadName                       string `mapstructure:"log_thread_name" cmdx:"publisher.kafka.client.log.thread.name" default:"true"`
	EnableRandomSeed                    string `mapstructure:"enable_random_seed" cmdx:"publisher.kafka.client.enable.random.seed" default:"true"`
	LogConnectionClose                  string `mapstructure:"log_connection_close" cmdx:"publisher.kafka.client.log.connection.close" default:"true"`
	InternalTerminationSignal           string `mapstructure:"internal_termination_signal" cmdx:"publisher.kafka.client.internal.termination.signal" default:"0"`
	ApiVersionRequest                   string `mapstructure:"api_version_request" cmdx:"publisher.kafka.client.api.version.request" default:"true"`
	ApiVersionRequestTimeoutMS          string `mapstructure:"api_version_request_timeout_ms" cmdx:"publisher.kafka.client.api.version.request.timeout.ms" default:"10000"`
	ApiVersionFallbackMS                string `mapstructure:"api_version_fallback_ms" cmdx:"publisher.kafka.client.api.version.fallback.ms" default:"0"`
	BrokerVersionFallback               string `mapstructure:"broker_version_fallback" cmdx:"publisher.kafka.client.broker.version.fallback" default:"0.10.0"`
	SecurityProtocol                    string `mapstructure:"security_protocol" cmdx:"publisher.kafka.client.security.protocol" default:"plaintext"`
	SSLCipherSuites                     string `mapstructure:"ssl_cipher_suites" cmdx:"publisher.kafka.client.ssl.cipher.suites" default:""`
	SSLCurvesList                       string `mapstructure:"ssl_curves_list" cmdx:"publisher.kafka.client.ssl.curves.list" default:""`
	SSLSigalgsList                      string `mapstructure:"ssl_sigalgs_list" cmdx:"publisher.kafka.client.ssl.sigalgs.list" default:""`
	SSLKeyLocation                      string `mapstructure:"ssl_key_location" cmdx:"publisher.kafka.client.ssl.key.location" default:""`
	SSLKeyPassword                      string `mapstructure:"ssl_key_password" cmdx:"publisher.kafka.client.ssl.key.password" default:""`
	SSLKeyPem                           string `mapstructure:"ssl_key_pem" cmdx:"publisher.kafka.client.ssl.key.pem" default:""`
	SSLCertificateLocation              string `mapstructure:"ssl_certificate_location" cmdx:"publisher.kafka.client.ssl.certificate.location" default:""`
	SSLCertificatePem                   string `mapstructure:"ssl_certificate_pem" cmdx:"publisher.kafka.client.ssl.certificate.pem" default:""`
	SSLCaLocation                       string `mapstructure:"ssl_ca_location" cmdx:"publisher.kafka.client.ssl.ca.location" default:""`
	SSLCaPem                            string `mapstructure:"ssl_ca_pem" cmdx:"publisher.kafka.client.ssl.ca.pem" default:""`
	SSLCrlLocation                      string `mapstructure:"ssl_crl_location" cmdx:"publisher.kafka.client.ssl.crl.location" default:""`
	SSLKeystoreLocation                 string `mapstructure:"ssl_keystore_location" cmdx:"publisher.kafka.client.ssl.keystore.location" default:""`
	SSLKeystorePassword                 string `mapstructure:"ssl_keystore_password" cmdx:"publisher.kafka.client.ssl.keystore.password" default:""`
	SSLEngineLocation                   string `mapstructure:"ssl_engine_location" cmdx:"publisher.kafka.client.ssl.engine.location" default:""`
	SSLEngineID                         string `mapstructure:"ssl_engine_id" cmdx:"publisher.kafka.client.ssl.engine.id" default:"dynamic"`
	EnableSSLCertificateVerification    string `mapstructure:"enable_ssl_certificate_verification" cmdx:"publisher.kafka.client.enable.ssl.certificate.verification" default:"true"`
	SSLEndpointIdentificationAlgorithm  string `mapstructure:"ssl_endpoint_identification_algorithm" cmdx:"publisher.kafka.client.ssl.endpoint.identification.algorithm" default:"none"`
	SASLMechanisms                      string `mapstructure:"sasl_mechanisms" cmdx:"publisher.kafka.client.sasl.mechanisms" default:"GSSAPI"`
	SASLMechanism                       string `mapstructure:"sasl_mechanism" cmdx:"publisher.kafka.client.sasl.mechanism" default:"GSSAPI"`
	SASLKerberosServiceName             string `mapstructure:"sasl_kerberos_service_name" cmdx:"publisher.kafka.client.sasl.kerberos.service.name" default:"kafka"`
	SASLKerberosPrincipal               string `mapstructure:"sasl_kerberos_principal" cmdx:"publisher.kafka.client.sasl.kerberos.principal" default:"kafkaclient"`
	SASLKerberosKinitCmd                string `mapstructure:"sasl_kerberos_kinit_cmd" cmdx:"publisher.kafka.client.sasl.kerberos.kinit.cmd" default:"kinit -R -t \"%{sasl.kerberos.keytab}\" -k \"%{sasl.kerberos.principal}\""`
	SASLKerberosKeytab                  string `mapstructure:"sasl_kerberos_keytab" cmdx:"publisher.kafka.client.sasl.kerberos.keytab" default:""`
	SASLKerberosMinTimeBeforeRelogin    string `mapstructure:"sasl_kerberos_min_time_before_relogin" cmdx:"publisher.kafka.client.sasl.kerberos.min.time.before.relogin" default:"60000"`
	SASLUsername                        string `mapstructure:"sasl_username" cmdx:"publisher.kafka.client.sasl.username" default:""`
	SASLPassword                        string `mapstructure:"sasl_password" cmdx:"publisher.kafka.client.sasl.password" default:""`
	SASLOauthbearerConfig               string `mapstructure:"sasl_oauthbearer_config" cmdx:"publisher.kafka.client.sasl.oauthbearer.config" default:""`
	EnableSASLOauthbearerUnsecureJwt    string `mapstructure:"enable_sasl_oauthbearer_unsecure_jwt" cmdx:"publisher.kafka.client.enable.sasl.oauthbearer.unsecure.jwt" default:"false"`
	SASLOauthbearerMethod               string `mapstructure:"sasl_oauthbearer_method" cmdx:"publisher.kafka.client.sasl.oauthbearer.method" default:"default"`
	SASLOauthbearerClientID             string `mapstructure:"sasl_oauthbearer_client_id" cmdx:"publisher.kafka.client.sasl.oauthbearer.client.id" default:""`
	SASLOauthbearerClientSecret         string `mapstructure:"sasl_oauthbearer_client_secret" cmdx:"publisher.kafka.client.sasl.oauthbearer.client.secret" default:""`
	SASLOauthbearerScope                string `mapstructure:"sasl_oauthbearer_scope" cmdx:"publisher.kafka.client.sasl.oauthbearer.scope" default:""`
	SASLOauthbearerExtensions           string `mapstructure:"sasl_oauthbearer_extensions" cmdx:"publisher.kafka.client.sasl.oauthbearer.extensions" default:""`
	SASLOauthbearerTokenEndpointUrl     string `mapstructure:"sasl_oauthbearer_token_endpoint_url" cmdx:"publisher.kafka.client.sasl.oauthbearer.token.endpoint.url" default:""`
	PluginLibraryPaths                  string `mapstructure:"plugin_library_paths" cmdx:"publisher.kafka.client.plugin.library.paths" default:""`
	ClientRack                          string `mapstructure:"client_rack" cmdx:"publisher.kafka.client.client.rack" default:""`
	TransactionalID                     string `mapstructure:"transactional_id" cmdx:"publisher.kafka.client.transactional.id" default:""`
	TransactionTimeoutMS                string `mapstructure:"transaction_timeout_ms" cmdx:"publisher.kafka.client.transaction.timeout.ms" default:"60000"`
	EnableIdempotence                   string `mapstructure:"enable_idempotence" cmdx:"publisher.kafka.client.enable.idempotence" default:"false"`
	EnableGaplessGuarantee              string `mapstructure:"enable_gapless_guarantee" cmdx:"publisher.kafka.client.enable.gapless.guarantee" default:"false"`
	QueueBufferingMaxMessages           string `mapstructure:"queue_buffering_max_messages" cmdx:"publisher.kafka.client.queue.buffering.max.messages" default:"100000"`
	QueueBufferingMaxKbytes             string `mapstructure:"queue_buffering_max_kbytes" cmdx:"publisher.kafka.client.queue.buffering.max.kbytes" default:"1048576"`
	QueueBufferingMaxMS                 string `mapstructure:"queue_buffering_max_ms" cmdx:"publisher.kafka.client.queue.buffering.max.ms" default:"5"`
	LingerMS                            string `mapstructure:"linger_ms" cmdx:"publisher.kafka.client.linger.ms" default:"5"`
	MessageSendMaxRetries               string `mapstructure:"message_send_max_retries" cmdx:"publisher.kafka.client.message.send.max.retries" default:"2147483647"`
	Retries                             string `mapstructure:"retries" cmdx:"publisher.kafka.client.retries" default:"2147483647"`
	RetryBackoffMS                      string `mapstructure:"retry_backoff_ms" cmdx:"publisher.kafka.client.retry.backoff.ms" default:"100"`
	QueueBufferingBackpressureThreshold string `mapstructure:"queue_buffering_backpressure_threshold" cmdx:"publisher.kafka.client.queue.buffering.backpressure.threshold" default:"1"`
	BatchNumMessages                    string `mapstructure:"batch_num_messages" cmdx:"publisher.kafka.client.batch.num.messages" default:"10000"`
	BatchSize                           string `mapstructure:"batch_size" cmdx:"publisher.kafka.client.batch.size" default:"1000000"`
	DeliveryReportOnlyError             string `mapstructure:"delivery_report_only_error" cmdx:"publisher.kafka.client.delivery.report.only.error" default:"false"`
	StickyPartitioningLingerMS          string `mapstructure:"sticky_partitioning_linger_ms" cmdx:"publisher.kafka.client.sticky.partitioning.linger.ms" default:"10"`
	RequestRequiredAcks                 string `mapstructure:"request_required_acks" cmdx:"publisher.kafka.client.request.required.acks" default:"-1"`
	Acks                                string `mapstructure:"acks" cmdx:"publisher.kafka.client.acks" default:"-1"`
	RequestTimeoutMS                    string `mapstructure:"request_timeout_ms" cmdx:"publisher.kafka.client.request.timeout.ms" default:"30000"`
	MessageTimeoutMS                    string `mapstructure:"message_timeout_ms" cmdx:"publisher.kafka.client.message.timeout.ms" default:"300000"`
	DeliveryTimeoutMS                   string `mapstructure:"delivery_timeout_ms" cmdx:"publisher.kafka.client.delivery.timeout.ms" default:"300000"`
	QueuingStrategy                     string `mapstructure:"queuing_strategy" cmdx:"publisher.kafka.client.queuing.strategy" default:"fifo"`
	Partitioner                         string `mapstructure:"partitioner" cmdx:"publisher.kafka.client.partitioner" default:"consistent_random"`
	CompressionCodec                    string `mapstructure:"compression_codec" cmdx:"publisher.kafka.client.compression.codec" default:"none"`
	CompressionType                     string `mapstructure:"compression_type" cmdx:"publisher.kafka.client.compression.type" default:"none"`
	CompressionLevel                    string `mapstructure:"compression_level" cmdx:"publisher.kafka.client.compression.level" default:"-1"`
}

type publisherKafka struct {
	FlushInterval       int               `mapstructure:"flush_interval_ms" cmdx:"publisher.kafka.flush.interval.ms" default:"1000"`
	DeliveryChannelSize int               `mapstructure:"delivery_channel_size" cmdx:"publisher.kafka.delivery.channel.size" default:"10"`
	ClientConfig        kafkaClientConfig `mapstructure:"client"`
}

func (k publisherKafka) ToKafkaConfigMap() *confluent.ConfigMap {
	configMap := &confluent.ConfigMap{}
	cfg := reflect.ValueOf(k.ClientConfig)
	for i := range cfg.NumField() {
		key := cfg.Type().Field(i).Tag.Get("mapstructure")
		key = strings.ReplaceAll(key, "_", ".")
		value := cfg.Field(i).String()
		configMap.SetKey(key, value)

	}
	return configMap
}

type publisher struct {
	Type    string           `mapstructure:"type" cmdx:"publisher.type" default:"kafka"`
	Kafka   publisherKafka   `mapstructure:"kafka"`
	PubSub  publisherPubSub  `mapstructure:"pubsub"`
	Kinesis publisherKinesis `mapstructure:"kinesis"`
}
