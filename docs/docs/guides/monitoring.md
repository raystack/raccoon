# Monitoring

Raccoon provides monitoring for server connection, Kafka publisher, resource usage, and event delivery. Reference for available metrics is available [here](../reference/metrics.md). The metrics are reported using [Statsd](https://www.datadoghq.com/blog/statsd/) protocol.

## How To Setup

```text
TL;DR
- Run Statsd supported metric collector
- Configure `METRIC_STATSD_ADDRESS` on Raccoon to send to the metric collector
- Visualize and create alerting from the collected metrics
```

Generally, you can follow the steps above and use any metric collector that supports statsd like [Telegraf](https://www.influxdata.com/blog/getting-started-with-sending-statsd-metrics-to-telegraf-influxdb/) or [Datadog](https://docs.datadoghq.com/developers/dogstatsd/?tab=hostagent).

This section will cover a setup example using [Telegraf](https://www.influxdata.com/time-series-platform/telegraf/), [Influx](https://www.influxdata.com/), [Kapacitor](https://www.influxdata.com/time-series-platform/kapacitor/), and [Grafana](https://grafana.com/) stack based on the steps above.

**Run Statsd Supported Metric Collector** To enable statsd on Telegraf you need to enable statsd input on `telegraf.conf` file. Following are default configurations that you can add based on statsd input [readme](https://github.com/influxdata/telegraf/blob/master/plugins/inputs/statsd/README.md).

```text
[[inputs.statsd]]
  protocol = "udp"
  max_tcp_connections = 250
  tcp_keep_alive = false
  service_address = ":8125"

  delete_gauges = true
  delete_counters = true
  delete_sets = true
  delete_timings = true

  percentiles = [50.0, 90.0, 99.0, 99.9, 99.95, 100.0]

  metric_separator = "_"

  parse_data_dog_tags = false
  datadog_extensions = false
  datadog_distributions = false

  allowed_pending_messages = 10000
  percentile_limit = 1000

[[outputs.influxdb]]
  urls = ["http://127.0.0.1:8086"]
  database = "raccoon"
  retention_policy = "autogen"
  write_consistency = "any"
  timeout = "5s"
```

**Configure Raccoon To Send To The Metric Collector** After you have the collector with the port configured, you need to set [METRIC\_STATSD\_ADDRESS](https://goto.gitbook.io/raccoon/reference/configurations#metric_statsd_address) to match the metric collector address. Suppose you deploy the telegraf using the default configuration above as sidecar or in localhost, you need to set the value to `:8125`.

**Visualize And Create Alerting From The Collected Metrics** Now that you have Raccoon and Telegraf as metric collector set, next is to use the metrics reported. You may notice that the Telegraf config above contains `outputs.influxdb`. That config will send the metric received to Influxdb. Make sure you have influx service accessible from the configured URL. You can visualize the metrics using Grafana. To do that, you need to [add influx datasource](https://www.influxdata.com/blog/how-grafana-dashboard-influxdb-flux-influxql/) to make the data available on Grafana. After that, you can use the data to You can visualize the metrics using Grafana. To do that, you need to [add influx datasource](https://www.influxdata.com/blog/how-grafana-dashboard-influxdb-flux-influxql/) to make the data available on Grafana. After that, you can use the data to [make dashboard](https://grafana.com/docs/grafana/latest/datasources/influxdb/#influxql-query-editor).

## Metrics Usages

Following are key monitoring statistics that you can infer from Raccoon metrics. Refer to this section to understand how to build alerting, dashboard, or analyze the metrics.

### Data Loss

To infer data loss, you can count [`kafka_messages_delivered_total`](https://goto.gitbook.io/raccoon/reference/metrics#kafka_messages_delivered_total) with tag `success=false`. You can also infer the loss rate by calculating the following.

`count(kafka_messages_delivered_total success=false)/count(kafka_messages_delivered_total)`

### Latency

Raccoon provides fine-grained metrics that denote latency. That gives clues as to where to look in case something goes wrong. In summary, these are key metrics for latency:

* [`event_processing_duration_milliseconds`](https://goto.gitbook.io/raccoon/reference/metrics#event_processing_duration_milliseconds) This metrics denotes overall latency. You need to look at other latency metrics to find the root cause when this metric is high.
* [`server_processing_latency_milliseconds`](https://goto.gitbook.io/raccoon/reference/metrics#server_processing_latency_milliseconds) Correlate this metric with `event_processing_duration_milliseconds` to infer whether the issue is with Raccoon itself, or something wrong with the network, or the way [sent\_time](https://github.com/goto/proton/blob/main/goto/raccoon/v1beta1/raccoon.proto#L47) is generated.-
* [`worker_processing_duration_milliseconds`](https://goto.gitbook.io/raccoon/reference/metrics#worker_processing_duration_milliseconds) High value of this metric indicates that the publisher is slow or can't keep up.

