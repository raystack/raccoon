| name | description | type | unit | source |
| ---- | ----- | ---- | --- | --- | 
| kafka_messages_delivered_total | Total messages delivered to kafka | counter | server.go |
| events_rx_bytes_total | Total bytes recieved from clients | counter | handler.go |
| events_rx_total | Total events recieved | counter | handler.go |
| ack_event_rtt_ms | Ack event rtt | histogram | ack.go |
| event_rtt_ms | Event Time RTT | histogram | ack.go |
| user_session_duration_milliseconds | Time user session is active | histogram | conn.go |
| batch_idle_in_channel_milliseconds |  | histogram | worker.go | 
| kafka_producebulk_tt_ms |  | histogram | worker.go | 
| event_processing_duration_milliseconds |  | histogram | worker.go | 
| worker_processing_duration_milliseconds |  | histogram | worker.go | 
| server_processing_latency_milliseconds |  | histogram | worker.go | 
| kafka_messages_delivered_total | | counter | kafka.go |
| kafka_unknown_topic_failure_total | | counter | kafka.go |
| batches_read_total | | counter | handler.go |
| events_rx_total | | counter | handler.go |
| events_duplicate_total | | counter | handler.go |
| server_ping_failure_total | | counter | pinger.go |
| conn_close_err_count | | counter | conn.go |
| user_connection_failure_total | | counter | upgrader.go |
| user_connection_success_total | | counter | upgrader.go |
| server_pong_failure_total | | counter | upgrader.go |
| server_go_routines_count_current | | gauge | server.go |
| server_mem_heap_alloc_bytes_current | | gauge | server.go |
| server_mem_heap_inuse_bytes_current | | gauge | server.go |
| server_mem_heap_objects_total_current | | gauge | server.go |
| server_mem_stack_inuse_bytes_current | | gauge | server.go |
| server_mem_gc_triggered_current | | gauge | server.go |
| server_mem_gc_pauseNs_current | | gauge | server.go |
| server_mem_gc_count_current | | gauge | server.go |
| server_mem_gc_pauseTotalNs_current | | gauge | server.go |
| kafka_tx_messages_total | | gauge | kafka.go |
| kafka_tx_messages_bytes_total | | gauge | kafka.go |
| kafka_brokers_tx_total | | gauge | kafka.go |
| kafka_brokers_tx_bytes_total | | gauge | kafka.go |
| kafka_brokers_rtt_average_milliseconds | | gauge | kafka.go |
| connections_count_current |  | gauge | kafka.go |