"use strict";(self.webpackChunkraccoon=self.webpackChunkraccoon||[]).push([[669],{5680:(e,n,t)=>{t.d(n,{xA:()=>p,yg:()=>d});var a=t(6540);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function r(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function s(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?r(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):r(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,a,o=function(e,n){if(null==e)return{};var t,a,o={},r=Object.keys(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(a=0;a<r.length;a++)t=r[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var i=a.createContext({}),c=function(e){var n=a.useContext(i),t=n;return e&&(t="function"==typeof e?e(n):s(s({},n),e)),t},p=function(e){var n=c(e.components);return a.createElement(i.Provider,{value:n},e.children)},u="mdxType",m={inlineCode:"code",wrapper:function(e){var n=e.children;return a.createElement(a.Fragment,{},n)}},g=a.forwardRef((function(e,n){var t=e.components,o=e.mdxType,r=e.originalType,i=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),u=c(t),g=o,d=u["".concat(i,".").concat(g)]||u[g]||m[g]||r;return t?a.createElement(d,s(s({ref:n},p),{},{components:t})):a.createElement(d,s({ref:n},p))}));function d(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var r=t.length,s=new Array(r);s[0]=g;var l={};for(var i in n)hasOwnProperty.call(n,i)&&(l[i]=n[i]);l.originalType=e,l[u]="string"==typeof e?e:o,s[1]=l;for(var c=2;c<r;c++)s[c]=t[c];return a.createElement.apply(null,s)}return a.createElement.apply(null,t)}g.displayName="MDXCreateElement"},9365:(e,n,t)=>{t.d(n,{A:()=>s});var a=t(6540),o=t(53);const r={tabItem:"tabItem_Ymn6"};function s(e){let{children:n,hidden:t,className:s}=e;return a.createElement("div",{role:"tabpanel",className:(0,o.A)(r.tabItem,s),hidden:t},n)}},4865:(e,n,t)=>{t.d(n,{A:()=>m});var a=t(8168),o=t(6540),r=t(53),s=t(2303),l=t(1682),i=t(4595),c=t(3104);const p={tabList:"tabList__CuJ",tabItem:"tabItem_LNqP"};function u(e){var n;const{lazy:t,block:s,defaultValue:u,values:m,groupId:g,className:d}=e,y=o.Children.map(e.children,(e=>{if((0,o.isValidElement)(e)&&"value"in e.props)return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)})),h=m??y.map((e=>{let{props:{value:n,label:t,attributes:a}}=e;return{value:n,label:t,attributes:a}})),f=(0,l.X)(h,((e,n)=>e.value===n.value));if(f.length>0)throw new Error(`Docusaurus error: Duplicate values "${f.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`);const b=null===u?u:u??(null==(n=y.find((e=>e.props.default)))?void 0:n.props.value)??y[0].props.value;if(null!==b&&!h.some((e=>e.value===b)))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${b}" but none of its children has the corresponding value. Available values are: ${h.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);const{tabGroupChoices:v,setTabGroupChoices:_}=(0,i.x)(),[N,k]=(0,o.useState)(b),T=[],{blockElementScrollPositionUntilNextRender:E}=(0,c.a_)();if(null!=g){const e=v[g];null!=e&&e!==N&&h.some((n=>n.value===e))&&k(e)}const w=e=>{const n=e.currentTarget,t=T.indexOf(n),a=h[t].value;a!==N&&(E(n),k(a),null!=g&&_(g,String(a)))},R=e=>{var n;let t=null;switch(e.key){case"ArrowRight":{const n=T.indexOf(e.currentTarget)+1;t=T[n]??T[0];break}case"ArrowLeft":{const n=T.indexOf(e.currentTarget)-1;t=T[n]??T[T.length-1];break}}null==(n=t)||n.focus()};return o.createElement("div",{className:(0,r.A)("tabs-container",p.tabList)},o.createElement("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,r.A)("tabs",{"tabs--block":s},d)},h.map((e=>{let{value:n,label:t,attributes:s}=e;return o.createElement("li",(0,a.A)({role:"tab",tabIndex:N===n?0:-1,"aria-selected":N===n,key:n,ref:e=>T.push(e),onKeyDown:R,onFocus:w,onClick:w},s,{className:(0,r.A)("tabs__item",p.tabItem,null==s?void 0:s.className,{"tabs__item--active":N===n})}),t??n)}))),t?(0,o.cloneElement)(y.filter((e=>e.props.value===N))[0],{className:"margin-top--md"}):o.createElement("div",{className:"margin-top--md"},y.map(((e,n)=>(0,o.cloneElement)(e,{key:n,hidden:e.props.value!==N})))))}function m(e){const n=(0,s.A)();return o.createElement(u,(0,a.A)({key:String(n)},e))}},2163:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>p,contentTitle:()=>i,default:()=>d,frontMatter:()=>l,metadata:()=>c,toc:()=>u});var a=t(8168),o=(t(6540),t(5680)),r=t(4865),s=t(9365);const l={},i="Deployment",c={unversionedId:"guides/deployment",id:"guides/deployment",title:"Deployment",description:"This section contains guides and suggestions related to Raccoon deployment.",source:"@site/docs/guides/deployment.md",sourceDirName:"guides",slug:"/guides/deployment",permalink:"/raccoon/guides/deployment",draft:!1,editUrl:"https://github.com/raystack/raccoon/edit/master/docs/docs/guides/deployment.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Publishing Events",permalink:"/raccoon/guides/publishing"},next:{title:"Monitoring",permalink:"/raccoon/guides/monitoring"}},p={},u=[{value:"Kubernetes",id:"kubernetes",level:2},{value:"Manifest",id:"manifest",level:3},{value:"Helm",id:"helm",level:3},{value:"Production Checklist",id:"production-checklist",level:2},{value:"Key Configurations",id:"key-configurations",level:3}],m={toc:u},g="wrapper";function d(e){let{components:n,...t}=e;return(0,o.yg)(g,(0,a.A)({},m,t,{components:n,mdxType:"MDXLayout"}),(0,o.yg)("h1",{id:"deployment"},"Deployment"),(0,o.yg)("p",null,"This section contains guides and suggestions related to Raccoon deployment."),(0,o.yg)("h2",{id:"kubernetes"},"Kubernetes"),(0,o.yg)("p",null,"Using ",(0,o.yg)("a",{parentName:"p",href:"https://hub.docker.com/r/raystack/raccoon"},"Raccoon docker image"),", you can deploy Raccoon on ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/"},"Kubernetes")," by specifying the image on the ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#creating-a-deployment"},"manifest"),". We also provide ",(0,o.yg)("a",{parentName:"p",href:"https://github.com/raystack/charts/tree/main/stable/raccoon"},"Helm chart")," to ease Kubernetes deployment. In this section we will cover simple deployment on Kubernetes using manifest and Helm."),(0,o.yg)("h3",{id:"manifest"},"Manifest"),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Prerequisite")),(0,o.yg)("ul",null,(0,o.yg)("li",{parentName:"ul"},"Kubernetes cluster setup"),(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("a",{parentName:"li",href:"https://kubernetes.io/docs/tasks/tools/#kubectl"},"Kubectl")," installed")),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Creating Kubernetes Resources")," You need at least 2 manifests for Raccoon. For ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/workloads/controllers/deployment"},"deployment")," and for ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/"},"configmap"),". Prepare both manifest as YAML file. You can fill in the configuration as needed."),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="configmap.yml"',title:'"configmap.yml"'},'apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: raccoon-config\n  namespace: default\n  labels:\n    application: raccoon\ndata:\n  PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS: "host.docker.internal:9093"\n  SERVER_WEBSOCKET_CONN_ID_HEADER: "X-User-ID"\n  SERVER_WEBSOCKET_PORT: "8080"\n\n  # depending on what monitoring stack you wish to use\n  # you can remove the statsd or prometheus config below\n  METRIC_STATSD_ENABLED: true\n  METRIC_STATSD_ADDRESS: "host.docker.internal:8125"\n  METRIC_PROMETHEUS_ENABLED: true\n  METRIC_PROMETHEUS_PORT: "9090"\n\n')),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="deployment.yaml"',title:'"deployment.yaml"'},"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: raccoon\n  labels:\n    application: raccoon\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      application: raccoon\n  template:\n    metadata:\n      labels:\n        application: raccoon\n      annotations:\n        # these are only necessary if you plan to use prometheus\n        # to collect metrics from raccoon. See \"Setting up monitoring\" below\n        # for more information.\n        prometheus.io/scrape: 'true'\n        prometheus.io/path: 'metrics'\n        prometheus.io/port: '9090'\n    spec:\n      containers:\n        - name: raccoon\n          image: \"raystack/raccoon:latest\"\n          imagePullPolicy: IfNotPresent\n          resources:\n            limits:\n              cpu: 200m\n              memory: 512Mi\n            requests:\n              cpu: 200m\n              memory: 512Mi\n          envFrom:\n            - configMapRef:\n                name: raccoon-config\n          env:\n            - name: POD_NAME\n              valueFrom:\n                fieldRef:\n                  fieldPath: metadata.name\n            - name: NODE_NAME\n              valueFrom:\n                fieldRef:\n                  fieldPath: spec.nodeName\n\n")),(0,o.yg)("p",null,"Suppose you save them as ",(0,o.yg)("inlineCode",{parentName:"p"},"configmap.yaml")," and ",(0,o.yg)("inlineCode",{parentName:"p"},"deployment.yaml"),". The next step is to apply the manifests to the Kubernetes cluster using ",(0,o.yg)("inlineCode",{parentName:"p"},"kubectl")," command."),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-bash"},"$ kubectl apply -f configmap.yaml -f deployment.yaml\n")),(0,o.yg)("p",null,"You'll find the resources are created. To see the status of the deployment, you can run following commands."),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-bash"},"# Check deployment status\n$ kubectl get deployment raccoon\n# Check configmap status\n$ kubectl get configmap raccoon-config\n")),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Configuration")," You can add or modify the configurations inside ",(0,o.yg)("inlineCode",{parentName:"p"},"configmap.yaml")," above. However, when you change the configmap, you also need to restart the deployment."),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Exposing Raccoon")," To make Raccoon accessible to the public, you need to setup the Kubernetes ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/services-networking/service/"},"service")," and ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/services-networking/ingress/"},"ingress"),". This setup may vary according to your need. There is plenty ",(0,o.yg)("a",{parentName:"p",href:"https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/"},"ingress controller")," you can choose. But first, you need to make sure that Websocket works with your choice of ingress controller."),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Setting up monitoring")),(0,o.yg)("p",null,"Raccoon supports ",(0,o.yg)("a",{parentName:"p",href:"https://github.com/statsd/statsd"},"statsd")," and ",(0,o.yg)("a",{parentName:"p",href:"https://prometheus.io/"},"prometheus")," as monitoring backends. The following section will guide on how to setup each of these."),(0,o.yg)(r.A,{mdxType:"Tabs"},(0,o.yg)(s.A,{value:"statsd",mdxType:"TabItem"},(0,o.yg)("p",null,"The recommended way of interfacing with statsd is to use ",(0,o.yg)("a",{parentName:"p",href:"https://www.influxdata.com/time-series-platform/telegraf/"},"telegraf"),".\ntelegraf is the open source server agent to help you collect metrics from your applications and push them to different data sources."),(0,o.yg)("p",null,"There are 2 options to integrate with Telegraf. One is to have Telegraf as a separate service and another is to have Telegraf as a sidecar. To have telegraf as a sidecar, you only need to add another configmap and another Telegraf container on the deployment above. You can add the container under ",(0,o.yg)("inlineCode",{parentName:"p"},"spec.template.spec.containers"),". Then, you can use default ",(0,o.yg)("inlineCode",{parentName:"p"},"METRIC_STATSD_ADDRESS")," which is ",(0,o.yg)("inlineCode",{parentName:"p"},":8125"),". Following is an example of Telegraf manifests that push to Influxdb."),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="deployment.yaml"',title:'"deployment.yaml"'},"\n---\ncontainers:\n  - image: telegraf:1.7.4-alpine\n    imagePullPolicy: IfNotPresent\n    name: telegrafd\n    resources:\n      limits:\n        cpu: 50m\n        memory: 64Mi\n      requests:\n        cpu: 50m\n        memory: 64Mi\n    volumeMounts:\n      - mountPath: /etc/telegraf\n        name: telegraf-conf\nvolumes:\n  - configMap:\n      defaultMode: 420\n      name: test-raccoon-telegraf-config\n    name: telegraf-conf\n")),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="telegraf-conf.yaml"',title:'"telegraf-conf.yaml"'},'apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: telegraf-conf\n  namespace: default\ndata:\n  telegraf.conf: |\n    [global_tags]\n      app = "test-raccoon"\n    [agent]\n      collection_jitter = "0s"\n      debug = false\n      flush_interval = "10s"\n      flush_jitter = "0s"\n      interval = "10s"\n      logfile = ""\n      metric_batch_size = 1000\n      metric_buffer_limit = 10000\n      omit_hostname = false\n      precision = ""\n      quiet = false\n      round_interval = true\n    [[outputs.influxdb]]\n      urls = ["http://localhost:8086"]\n      database = "test-db"\n      retention_policy = "autogen"\n      write_consistency = "any"\n      timeout = "5s"\n    [[inputs.statsd]]\n      allowed_pending_messages = 10000\n      delete_counters = true\n      delete_gauges = true\n      delete_sets = true\n      delete_timings = true\n      metric_separator = "."\n      parse_data_dog_tags = true\n      percentile_limit = 1000\n      percentiles = [\n        50,\n        95,\n        99\n      ]\n      service_address = ":8125"\n'))),(0,o.yg)(s.A,{value:"prometheus",mdxType:"TabItem"},(0,o.yg)("p",null,"You can run prometheus as a seperate service on your kubernetes cluster and configure it to scrape the metrics from your raccoon deployment."),(0,o.yg)("p",null,"Here's an example setup:"),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="prometheus-conf.yaml"',title:'"prometheus-conf.yaml"'},"apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: prometheus-config\n  labels:\n    app: prometheus\ndata:\n  prometheus.yml: |\n    global:\n      scrape_interval: 15s\n      evaluation_interval: 15s\n    scrape_configs:\n      - job_name: 'kubernetes-pods'\n        kubernetes_sd_configs:\n        - role: pod\n        relabel_configs:\n        - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]\n          action: keep\n          regex: true\n        - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]\n          action: replace\n          target_label: __metrics_path__\n          regex: (.+)\n          replacement: ${1}\n        - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]\n          action: replace\n          target_label: __address__\n          regex: (.+):(?:\\d+);(\\d+)\n          replacement: ${1}:${2}\n\n")),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="prometheus-deployment.yaml"',title:'"prometheus-deployment.yaml"'},"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: prometheus\n  labels:\n    app: prometheus\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: prometheus\n  template:\n    metadata:\n      labels:\n        app: prometheus\n    spec:\n      serviceAccountName: prometheus\n      containers:\n      - name: prometheus\n        image: prom/prometheus:latest\n        args:\n          - --config.file=/etc/prometheus/prometheus.yml\n          - --storage.tsdb.path=/prometheus/\n        ports:\n        - containerPort: 9090\n        volumeMounts:\n        - name: prometheus-config-volume\n          mountPath: /etc/prometheus/\n      volumes:\n      - name: prometheus-config-volume\n        configMap:\n          name: prometheus-config\n\n")),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="prometheus-service-account.yml"',title:'"prometheus-service-account.yml"'},"apiVersion: v1\nkind: ServiceAccount\nmetadata:\n  name: prometheus\n  labels:\n    app: prometheus\n\n")),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="prometheus-cluster-role.yml"',title:'"prometheus-cluster-role.yml"'},'apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRole\nmetadata:\n  name: prometheus\nrules:\n- apiGroups: [""]\n  resources:\n  - nodes\n  - nodes/proxy\n  - services\n  - endpoints\n  - pods\n  verbs: ["get", "list", "watch"]\n\n')),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-yaml",metastring:'title="prometheus-cluster-role-binding.yml"',title:'"prometheus-cluster-role-binding.yml"'},"apiVersion: rbac.authorization.k8s.io/v1\nkind: ClusterRoleBinding\nmetadata:\n  name: prometheus\nroleRef:\n  apiGroup: rbac.authorization.k8s.io\n  kind: ClusterRole\n  name: prometheus\nsubjects:\n- kind: ServiceAccount\n  name: prometheus\n  namespace: default\n\n")),(0,o.yg)("p",null,"Create these files on your local system that's configured to talk to your Kubernetes cluster, then deploy prometheus by running:"),(0,o.yg)("pre",null,(0,o.yg)("code",{parentName:"pre",className:"language-bash"},"$ kubectl apply -f prometheus-service-account.yml\n$ kubectl apply -f prometheus-cluster-role-binding.yml\n$ kubectl apply -f prometheus-cluster-role.yml\n$ kubectl apply -f prometheus-conf.yml\n$ kubectl apply -f prometheus-deployment.yml\n")),(0,o.yg)("p",null,"This setups uses prometheus's ",(0,o.yg)("a",{parentName:"p",href:"https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config"},"kubernetes_sd_config"),", which is a feature that allows prometheus to leverage service discovery to find the pods that you wish to scrape metrics from."))),(0,o.yg)("h3",{id:"helm"},"Helm"),(0,o.yg)("p",null,(0,o.yg)("strong",{parentName:"p"},"Prerequisite")),(0,o.yg)("ul",null,(0,o.yg)("li",{parentName:"ul"},"Kubernetes cluster setup"),(0,o.yg)("li",{parentName:"ul"},"Helm installed")),(0,o.yg)("p",null,"Raccoon has a Helm chart maintained on ",(0,o.yg)("a",{parentName:"p",href:"https://github.com/raystack/charts/tree/main/stable/raccoon"},"different repository"),". Refer to the repository for a complete deployment guide."),(0,o.yg)("h2",{id:"production-checklist"},"Production Checklist"),(0,o.yg)("p",null,"Before going to production, we recommend the following setup tips."),(0,o.yg)("h3",{id:"key-configurations"},"Key Configurations"),(0,o.yg)("p",null,"Followings are main configurations closely related to deployment that you need to pay attention:"),(0,o.yg)("ul",null,(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("p",{parentName:"li"},(0,o.yg)("a",{parentName:"p",href:"/raccoon/reference/configurations#server_websocket_port"},(0,o.yg)("inlineCode",{parentName:"a"},"SERVER_WEBSOCKET_PORT")))),(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("p",{parentName:"li"},(0,o.yg)("a",{parentName:"p",href:"/raccoon/reference/configurations#event_distribution_publisher_pattern"},(0,o.yg)("inlineCode",{parentName:"a"},"EVENT_DISTRIBUTION_PUBLISHER_PATTERN")))),(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("p",{parentName:"li"},(0,o.yg)("a",{parentName:"p",href:"/raccoon/reference/configurations#publisher_kafka_client_bootstrap_servers"},(0,o.yg)("inlineCode",{parentName:"a"},"PUBLISHER_KAFKA_CLIENT_BOOTSTRAP_SERVERS")))),(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("p",{parentName:"li"},(0,o.yg)("a",{parentName:"p",href:"/raccoon/reference/configurations#metric_statsd_address"},(0,o.yg)("inlineCode",{parentName:"a"},"METRIC_STATSD_ADDRESS")))),(0,o.yg)("li",{parentName:"ul"},(0,o.yg)("p",{parentName:"li"},(0,o.yg)("a",{parentName:"p",href:"/raccoon/reference/configurations#server_websocket_conn_id_header"},(0,o.yg)("inlineCode",{parentName:"a"},"SERVER_WEBSOCKET_CONN_ID_HEADER"))),(0,o.yg)("p",{parentName:"li"},(0,o.yg)("strong",{parentName:"p"},"TLS/HTTPS")),(0,o.yg)("p",{parentName:"li"},"Raccoon doesn't provide HTTPS on the application level. To enable HTTPS, you can maintain API gateway which terminates SSL connection. From API gateway onward, the connection is considered to be safe. For example, if you are deploying on Kubernetes, you can have an ingress setup and have SSL termination."),(0,o.yg)("p",{parentName:"li"},(0,o.yg)("strong",{parentName:"p"},"Authentication")),(0,o.yg)("p",{parentName:"li"},"Raccoon doesn't provide authentication on its own. However, you can still enable authentication by having it as a separate service. Then, you can use an API gateway to validate the authentication using a token."),(0,o.yg)("p",{parentName:"li"},(0,o.yg)("strong",{parentName:"p"},"Test The Setup")),(0,o.yg)("p",{parentName:"li"},"To make sure the deployment can handle the load, you need to test it with the same number of connections and request you are expecting. You can find a guide on how to publish events ",(0,o.yg)("a",{parentName:"p",href:"/raccoon/guides/publishing"},"here"),". You can also check client example ",(0,o.yg)("a",{parentName:"p",href:"/raccoon/quickstart##publishing-your-first-event"},"here"),". If there is something wrong with Raccoon, you can check the ",(0,o.yg)("a",{parentName:"p",href:"/raccoon/guides/troubleshooting"},"troubleshooting")," section."))))}d.isMDXComponent=!0}}]);