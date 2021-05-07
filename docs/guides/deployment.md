# Deployment
This section contains guides and advices related to Raccoon deployment.
## Standalone
This section contains standalone deployment options for Raccoon server. You can use this deployment method on your local machine to get started with Raccoon.
### Local Machine
**Prerequisite**
- GO 1.14 or higher installed
- Unix based machine

**Run The Executable**
Before you can run the server on your local machine. First, you need to compile it using the following command.
```sh
$ git clone https://github.com/odpf/raccoon
$ cd raccoon
$ make
```
You can find the executable on `./out/raccoon`.
```sh
# Run the executable
$ ./out/raccoon
```

**Configuration**
You can have the [configuration](https://odpf.gitbook.io/raccoon/reference/configurations) set either in `.env` file where you run the executable or export it as env variable.

### Docker
**Prerequisite**
- Docker installed

**Run Docker Image**
Raccoon provides Docker [image](https://hub.docker.com/r/odpf/raccoon) as part of the release. To run Raccoon with default port exposed on `localhost` you can run the following.
```sh
$ docker pull odpf/raccoon:latest
$ docker run -p 8080:8080 \
  -e SERVER_WEBSOCKET_CONN_UNIQ_ID_HEADER=x-user-id \
  -e PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS=host.docker.internal:9092 \
  -e METRIC_STATSD_ADDRESS=host.docker.internal:8125 \
  odpf/raccoon:latest
```

**Configuration**
- Use `-e KEY=VALUE` flag  when you do `docker run`
- Configurations related to URL assumes that the services run on `localhost`. When you run Raccoon inside docker, make sure to change the URL accordingly. For example, you can use docker special DNS name `host.docker.internal` to resolve to your host machine network.
## Docker Compose
**Prerequisite**
- Docker installed

**Running The Services**
This repository contains `docker-compose.yml` file for development and integration test purpose. The docker-compose deploys Kafka and Zookeeper along with the Raccoon service. One usecase for this setup is when you are developing client for Raccoon, and need to test that client.
You can run the docker compose from the make script.
```sh
# Build and up
$ make docker-run
# To stop the services without deleting the container
$ make docker-stop
# To continue the existing services. If you make changes to the codebase, you need to run `make docker-run` instead
$ make docker-start
```

You can consume the published events from the host machine by using `localhost:9094` as kafka broker server. Mind the [topic routing](https://odpf.gitbook.io/raccoon/concepts/architecture#event-distribution) when you consume the events.
**Configuration**
Since this setup is using local `Dockerfile`, you can provide the configuration as `.env` file. Before you run the docker compose, you need to set `PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS=kafka:9092`.
## Kubernetes
Using Raccoon docker image, you can deploy Raccoon on [Kubernetes](https://kubernetes.io/) by specifying the image on the [manifest](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment). We also provides [Helm chart](https://github.com/odpf/charts/tree/main/stable/raccoon) to ease Kubernetes deployment. This section we will cover simple deployment on Kubernetes using manifest and Helm.
### Manifest
**Prerequisite**
- Kubernetes cluster setup
- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) installed

**Creating Kubernetes Resources**
You need at least 2 mainfest for Raccoon. For [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment) and for [configmap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/). Prepare both manifest as yaml file. You can fill the configuration as needed.
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: raccoon-config
  namespace: default
  labels:
    application: raccoon
data:
  METRIC_STATSD_ADDRESS: "host.docker.internal:8125"
  PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS: "host.docker.internal:9093"
  SERVER_WEBSOCKET_CONN_UNIQ_ID_HEADER: "x-user-id"
  SERVER_WEBSOCKET_PORT: "8080"
```
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: raccoon
  labels:
    application: raccoon
spec:
  replicas: 1
  selector:
    matchLabels:
      application: raccoon
  template:
    metadata:
      labels:
        application: raccoon
    spec:
      containers:
      - name: raccoon
        image: "odpf/raccoon:latest"
        imagePullPolicy: IfNotPresent
        resources:
            limits:
              cpu: 200m
              memory: 512Mi
            requests:
              cpu: 200m
              memory: 512Mi
        envFrom:
          - configMapRef:
              name: raccoon-config
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
      volumes:
```

Suppose you save them as `configmap.yaml` and `deployment.yaml`. The next step is to apply the manifests to the kube cluster using `kubectl` command.
```sh
$ kubectl apply -f configmap.yaml -f deployment.yaml
```
You'll find the resources are created. To see the status of the deployment you can run
```sh
# Check deployment status
$ kubectl get deployment raccoon
# Check configmap status
$ kubectl get configmap raccoon-config
```

**Configuration**
You can add or modify the configurations inside `configmap.yaml` above. When you change the configmap, you also need to restart the deployment.

### Helm
**Prerequisite**
- Kubernetes cluster setup
- Helm installed

Raccoon has helm chart maintained on [different repository](https://github.com/odpf/charts/tree/main/stable/raccoon). Refer to the repository for complete deployment guide.

## Production Checklist
Before going to production, we recommend following setup tips.
### TLS/HTTPS
Raccoon doesn't provide HTTPS on the application level. To enable HTTPS, you can maintain API gateway which terminates SSL connection. From API gateway onward, the connection is considered to be safe. For example, if you are deploying on Kubernetes, you can have an ingress setup and have SSL termination.
### Authentication
Raccoon doesn't provide authentication on it's own. However, you can still enable authentication by having it as separate service. Then, you can use API gateway to validate the authentication using token.
### Test The Setup
To make sure the deployment can handle the load, you need to test it with the same number of connections and request you are expecting. You can find guide on how to publish events [here](https://odpf.gitbook.io/raccoon/guides/publishing). If there is something wrong with Raccon, you can check the [troubleshooting](https://odpf.gitbook.io/raccoon/guides/troubleshooting) section.
