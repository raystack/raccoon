# Deployment

This section contains guides and suggestions related to Raccoon deployment.

## Kubernetes

Using [Raccoon docker image](https://hub.docker.com/r/goto/raccoon), you can deploy Raccoon on [Kubernetes](https://kubernetes.io/) by specifying the image on the [manifest](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment). We also provide [Helm chart](https://github.com/goto/charts/tree/main/stable/raccoon) to ease Kubernetes deployment. In this section we will cover simple deployment on Kubernetes using manifest and Helm.

### Manifest

**Prerequisite**

* Kubernetes cluster setup
* [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) installed

**Creating Kubernetes Resources** You need at least 2 manifests for Raccoon. For [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment) and for [configmap](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/). Prepare both manifest as YAML file. You can fill in the configuration as needed.

`configmap.yaml`

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
  SERVER_WEBSOCKET_CONN_ID_HEADER: "X-User-ID"
  SERVER_WEBSOCKET_PORT: "8080"
```

`deployment.yaml`

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
        image: "gotocompany/raccoon:latest"
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

Suppose you save them as `configmap.yaml` and `deployment.yaml`. The next step is to apply the manifests to the Kubeenetes cluster using `kubectl` command.

```bash
$ kubectl apply -f configmap.yaml -f deployment.yaml
```

You'll find the resources are created. To see the status of the deployment, you can run following commands.

```bash
# Check deployment status
$ kubectl get deployment raccoon
# Check configmap status
$ kubectl get configmap raccoon-config
```

**Configuration** You can add or modify the configurations inside `configmap.yaml` above. However, when you change the configmap, you also need to restart the deployment.

**Exposing Raccoon** To make Raccoon accessible to the public, you need to setup the Kubernetes [service](https://kubernetes.io/docs/concepts/services-networking/service/) and [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/). This setup may vary according to your need. There is plenty [ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/) you can choose. But first, you need to make sure that Websocket works with your choice of ingress controller.

**Integrating With Telegraf** There are 2 options to integrate with Telegraf. One is to have Telegraf as separate service another is to have Telegraf as a sidecar. To have telegraf as a sidecar, you only need to add another configmap and another Telegraf container on the deployment above. You can add the container under `spec.template.spec.containers`. Then, you can use default `METRIC_STATSD_ADDRESS` which is `:8125`. Following is an example of Telegraf manifests that push to Influxdb.

`deployment.yaml`

```yaml
...
      containers:
      - image: telegraf:1.7.4-alpine
        imagePullPolicy: IfNotPresent
        name: telegrafd
        resources:
          limits:
            cpu: 50m
            memory: 64Mi
          requests:
            cpu: 50m
            memory: 64Mi
        volumeMounts:
        - mountPath: /etc/telegraf
          name: telegraf-conf
      volumes:
      - configMap:
          defaultMode: 420
          name: test-raccoon-telegraf-config
        name: telegraf-conf  
...
```

`telegraf-conf.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: telegraf-conf
  namespace: default
data:
  telegraf.conf: |
    [global_tags]
      app = "test-raccoon"
    [agent]
      collection_jitter = "0s"
      debug = false
      flush_interval = "10s"
      flush_jitter = "0s"
      interval = "10s"
      logfile = ""
      metric_batch_size = 1000
      metric_buffer_limit = 10000
      omit_hostname = false
      precision = ""
      quiet = false
      round_interval = true
    [[outputs.influxdb]]
      urls = ["http://localhost:8086"]
      database = "test-db"
      retention_policy = "autogen"
      write_consistency = "any"
      timeout = "5s"
    [[inputs.statsd]]
      allowed_pending_messages = 10000
      delete_counters = true
      delete_gauges = true
      delete_sets = true
      delete_timings = true
      metric_separator = "."
      parse_data_dog_tags = true
      percentile_limit = 1000
      percentiles = [
        50,
        95,
        99
      ]
      service_address = ":8125"
```

### Helm

**Prerequisite**

* Kubernetes cluster setup
* Helm installed

Raccoon has a Helm chart maintained on [different repository](https://github.com/goto/charts/tree/main/stable/raccoon). Refer to the repository for a complete deployment guide.

## Production Checklist

Before going to production, we recommend the following setup tips.

### Key Configurations

Followings are main configurations closely related to deployment that you need to pay attention:

* [`SERVER_WEBSOCKET_PORT`](https://goto.gitbook.io/raccoon/reference/configurations#server_websocket_port)
* [`EVENT_DISTRIBUTION_PUBLISHER_PATTERN`](https://goto.gitbook.io/raccoon/reference/configurations#event_distribution_publisher_pattern)
* [`PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS`](https://goto.gitbook.io/raccoon/reference/configurations#publisher_kafka_client_bootstrap_servers)
* [`METRIC_STATSD_ADDRESS`](https://goto.gitbook.io/raccoon/reference/configurations#metric_statsd_address)
* [`SERVER_WEBSOCKET_CONN_ID_HEADER`](https://goto.gitbook.io/raccoon/reference/configurations#server_websocket_conn_id_header)

  **TLS/HTTPS**

  Raccoon doesn't provide HTTPS on the application level. To enable HTTPS, you can maintain API gateway which terminates SSL connection. From API gateway onward, the connection is considered to be safe. For example, if you are deploying on Kubernetes, you can have an ingress setup and have SSL termination.

  **Authentication**

  Raccoon doesn't provide authentication on its own. However, you can still enable authentication by having it as a separate service. Then, you can use an API gateway to validate the authentication using a token.

  **Test The Setup**

  To make sure the deployment can handle the load, you need to test it with the same number of connections and request you are expecting. You can find a guide on how to publish events [here](https://goto.gitbook.io/raccoon/guides/publishing). You can also check example client [here](https://github.com/goto/raccoon/tree/main/docs/example). If there is something wrong with Raccon, you can check the [troubleshooting](https://goto.gitbook.io/raccoon/guides/troubleshooting) section.

