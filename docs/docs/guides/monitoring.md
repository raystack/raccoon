import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";

# Monitoring

Raccoon provides monitoring for server connection, publisher, resource usage, and event delivery. Reference for available metrics is available [here](reference/metrics.md). The metrics are reported using [Statsd](https://www.datadoghq.com/blog/statsd/) and [Prometheus](https://prometheus.io/)


## How To Setup

```mdx-code-block
<Tabs>
<TabItem value="statsd">
```

```text
TL;DR
- Run Statsd supported metric collector
- Configure `METRIC_STATSD_ADDRESS` on Raccoon to send to the metric collector
- Visualize and create alerting from the collected metrics
```

Generally, you can follow the steps above and use any metric collector that supports statsd like [Telegraf](https://www.influxdata.com/blog/getting-started-with-sending-statsd-metrics-to-telegraf-influxdb/) or [Datadog](https://docs.datadoghq.com/developers/dogstatsd/?tab=hostagent).

This section will cover a setup example using [Telegraf](https://www.influxdata.com/time-series-platform/telegraf/), [Influx](https://www.influxdata.com/), [Kapacitor](https://www.influxdata.com/time-series-platform/kapacitor/), and [Grafana](https://grafana.com/) stack based on the steps above.

**Run Statsd Supported Metric Collector** To enable statsd on Telegraf you need to enable statsd input on `telegraf.conf` file. Following are default configurations that you can add based on statsd input [readme](https://github.com/influxdata/telegraf/blob/master/plugins/inputs/statsd/README.md).

```toml
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

**Configure Raccoon To Send To The Metric Collector** After you have the collector with the port configured, you need to set [METRIC_STATSD_ADDRESS](reference/configurations.md#metric_statsd_address) to match the metric collector address. Suppose you deploy the telegraf using the default configuration above as sidecar or in localhost, you need to set the value to `:8125`.

**Visualize And Create Alerting From The Collected Metrics** Now that you have Raccoon and Telegraf as metric collector set, next is to use the metrics reported. You may notice that the Telegraf config above contains `outputs.influxdb`. That config will send the metric received to Influxdb. Make sure you have influx service accessible from the configured URL. You can visualize the metrics using Grafana. To do that, you need to [add influx datasource](https://www.influxdata.com/blog/how-grafana-dashboard-influxdb-flux-influxql/) to make the data available on Grafana. After that, you can use the data to You can [visualize the metrics](https://grafana.com/docs/grafana/latest/datasources/influxdb/#influxql-query-editor) using Grafana. 

```mdx-code-block
</TabItem>
<TabItem value="prometheus">
```
```text
TL;DR
- Run Prometheus
- Configure `METRIC_PROMETHEUS_PORT` and `METRIC_PROMETHEUS_ENABLED` on Raccoon
- Visualize and create alerting from the collected metrics
```
Setting up [Prometheus](https://prometheus.io) is fairly straight-forward. Prometheus is available as a self-contained binary program for most platforms.
For alerting you can use [alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) that let's you define alerts and offers integration with different notification platforms. For visualisation [Grafana](https://grafana.com/) comes with out-of-the-box support for prometheus as a data source.

You can download prometheus from their [official website](https://prometheus.io/download/).

**Run Prometheus**

Let's explore an example setup that runs prometheus locally. Begin by creating a new directory and downloading prometheus.
```bash
$ mkdir prometheus-for-raccoon
$ cd prometheus-for-raccoon
$ wget https://github.com/prometheus/prometheus/releases/download/v2.53.1/prometheus-2.53.1.linux-amd64.tar.gz
$ tar xzvf prometheus-2.53.1.linux-amd64.tar.gz
$ cd prometheus-2.53.1.linux-amd64
```

Next, we will edit `prometheus.yml` to tell prometheus to scrape metrics from raccoon. You can use any text editor that you're familiar with.


```yaml title=prometheus.yml
global:
  scrape_interval: 15s 
  evaluation_interval: 15s 

scrape_configs:
  - job_name: "raccoon"
    static_configs:
      - targets: ["localhost:8888"]

```

Now run prometheus 
```bash
$ ./prometheus --config.file=./prometheus.yml
```

We have now configured prometheus to scrape metrics from `localhost:8888`. We will now tell raccoon to expose prometheus metric on this port.

**Configure `METRIC_PROMETHEUS_PORT` and `METRIC_PROMETHEUS_ENABLED` on Raccoon**

By default, raccoon doesn't expose prometheus metrics. To enable prometheus metrics, you need to set the following environment variables:

```bash
METRIC_PROMETHEUS_ENABLED=true
METRIC_PROMETHEUS_PORT=8888   
```
Now when you run raccoon, prometheus will start collecting metrics from it.

**Visualize And Create Alerting From The Collected Metrics** 

Now that you have Raccoon and Prometheus setup, next is to use the metrics reported. You can visualize the metrics using Grafana. To do that, you need to [add prometheus datasource](https://grafana.com/docs/grafana/latest/datasources/prometheus/) to make the data available on Grafana. After that, you can use the data to visualize the metrics using Grafana. 


```mdx-code-block
</TabItem>
</Tabs>
```



## Metrics Usages

Following are key monitoring statistics that you can infer from Raccoon metrics. Refer to this section to understand how to build alerting, dashboard, or analyze the metrics.

### Data Loss

To infer data loss, you can count [`kafka_messages_delivered_total`](reference/metrics.md#kafka_messages_delivered_total) with tag `success=false`. You can also infer the loss rate by calculating the following.

`count(kafka_messages_delivered_total success=false)/count(kafka_messages_delivered_total)`

For other publishers, just replace `kafka` in the metric name with the name of the publisher. For instance, analogs of `kafka_messages_delivered_total` for PubSub and Kinesis would be:
* [`pubsub_messages_delivered_total`](reference/metrics.md#pubsub_messages_delivered_total)
* [`kinesis_messages_delivered_total`](reference/metrics.md#kinesis_messages_delivered_total)


### Latency

Raccoon provides fine-grained metrics that denote latency. That gives clues as to where to look in case something goes wrong. In summary, these are key metrics for latency:

- [`event_processing_duration_milliseconds`](reference/metrics.md#event_processing_duration_milliseconds) This metrics denotes overall latency. You need to look at other latency metrics to find the root cause when this metric is high.
- [`server_processing_latency_milliseconds`](reference/metrics.md#server_processing_latency_milliseconds) Correlate this metric with `event_processing_duration_milliseconds` to infer whether the issue is with Raccoon itself, or something wrong with the network, or the way [sent_time](https://github.com/raystack/proton/blob/main/raystack/raccoon/v1beta1/raccoon.proto#L47) is generated.-
- [`worker_processing_duration_milliseconds`](reference/metrics.md#worker_processing_duration_milliseconds) High value of this metric indicates that the publisher is slow or can't keep up.
