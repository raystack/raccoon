# Quickstart

This document will guide you on how to get Raccoon server running on your system. This document assumes that you have Docker installed on your system

## Run Raccoon with Docker 

Here's a minimal setup that runs Raccoon with `log` publisher. 

```bash
$ docker run --rm -p 8080:8080 \
  raystack/raccoon server \
  --publisher.type "log" \
  --server.websocket.conn.id.header "x-user-id"
```

To test whether Raccoon is running or not, you can try to ping the server.

```bash
$ curl http://localhost:8080/ping
```

## Publishing Your First Event

You can use `curl` to publish events to raccoon's REST API

```bash
$ curl -XPOST "http://localhost:8080/api/v1/events" \
    -H "content-type: application/json" \
    -H "X-User-ID: user-one" \
    -d "
{
    \"req_guid\": \"foobar-123\",
    \"sent_time\": {
        \"seconds\": $(date +%s),
        \"nanos\": $(date +%N)
    },
    \"events\": [
        {
            \"type\": \"page\",
            \"event_bytes\": \"$(echo \"EVENT\" | base64)\"
        }
    ]
}"
```

## Where To Go Next

* [Client Documentation](clients/overview.md)
* [Publishing Guide](guides/publishing.md)
* [Architecture](concepts/architecture.md)
