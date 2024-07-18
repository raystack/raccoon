# Benchmarks

This page contains performance benchmarks for raccoon.

## Table of Contents

* [Websocket](benchmarks.md#websocket)

## WebSocket

### Setup

Raccoon benchmarking was done using a client which creates multiple parallel connections to Raccoon in one go and then sends batches of events every 10 seconds.
This whole setup was deployed on a Kubernetes cluster running on GCP with one or multiple raccoon pods.

### Result

Following are the benchmarking results for various versions

| Raccoon version | Duration  | No. of Connections | No. of Raccoon Pods | No. of events/10s | Server Processing Latency(P95) | Server Processing Latency(Upper) | Workers Latency (mean p95)| Workers Latency (max upper) | Memory Used per pod |   CPU Cores Used per pod   |
|-----------------|-----------|--------------------|---------------------|-------------------|--------------------------------|----------------------------------|---------------------------|-----------------------------|---------------------|----------------------------|
| v0.1.0          | 1 hour    |       10000        |         1           |       1500000     |           ~ 6 - 22 ms          |           ~ 3 - 913 ms           |        ~ 1.7 - 7 ms       |       ~ 2 - 140 ms          |    ~ 711 - 870 MB   |           ~ 2.2            |
| v0.1.0          | 1 hour    |       50000        |         3           |       25000       |           ~ 20 - 30 ms         |              ~ 3 s               |        ~ 35 - 50 ms       |       ~ 35 - 45 ms          |    ~ 1.0- 1.5 GB    |         ~ 0.3 - 0.6        |
| v0.1.0          | 1 hour    |       50000        |         3           |       100000      |           ~ 3 - 5 ms           |              ~ 2.5 s             |        ~ 20-  30 ms       |       ~ 20 - 30 ms          |     400 - 500 MB    |            ~ 0.7           |
| v0.1.0          | 1 hour    |       100000       |         5           |       100000      |           ~ 3 - 9 ms           |              ~ 2.5 s             |        ~ 20-  30 ms       |       ~ 20 - 30 ms          |     ~ 1.7 - 2GB     |             ~ 1            |
| v0.1.2          | 30 min    |       10000        |         1           |       1500000     |           ~ 1 - 7.13k ms       |           ~ 1 - 7.13k ms         |      ~ 0.8 - 9.25 ms      |       ~ 2 - 224 ms          |     ~ 960MB -1.2GB  |          ~ ~ 2.57          |
