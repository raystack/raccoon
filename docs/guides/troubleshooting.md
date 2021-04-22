# Troubleshooting

## OS Configuration

Raccoon is using WebSocket as communication protocol from client to server. Websocket requires maintaining long-running connections. Each connection costs the OS an open file descriptor. When you reach the limit of the configured open file descriptor, the server won't be able to accept new connection. By default, OS limit the number of open file descriptor. You can look up how to increase the max open file descriptor. On Unix you can do `ulimit -n` to check max open file descriptor and `ulimit -n <number>` to set new limit.

## Tuning Raccoon

To figure out the right configuration for production, you need to load test the deployment with the same throughput on production. You can tune the configuration based on the load test result. Important configurations to tune are:

* [SERVER\_WEBSOCKET\_MAX\_CONN](https://github.com/odpf/raccoon/tree/d5c19b5cc8e75ffda16385e70d89e082146e044f/docs/reference/configuration.md#Server#SERVER_WEBSOCKET_MAX_CONN)
* [WORKER\_BUFFER\_CHANNEL\_SIZE](https://github.com/odpf/raccoon/tree/d5c19b5cc8e75ffda16385e70d89e082146e044f/docs/reference/configuration.md#Worker#WORKER_BUFFER_CHANNEL_SIZE)
* [WORKER\_POOL\_SIZE](https://github.com/odpf/raccoon/tree/d5c19b5cc8e75ffda16385e70d89e082146e044f/docs/reference/configuration.md#Worker#WORKER_POOL_SIZE)
* [PUBLISHER\_KAFKA\_CLIENT\_BATCH\_NUM\_MESSAGES](https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md)
* [PUBLISHER_KAFKA\_CLIENT_\*](https://github.com/odpf/raccoon/tree/d5c19b5cc8e75ffda16385e70d89e082146e044f/docs/reference/configuration.md#Publisher#PUBLISHER_KAFKA_CLIENT_*)

