package linkerd2

import (
	"github.com/ghodss/yaml"
	"github.com/iancoleman/strcase"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := writeTemplates("hack/install/linkerd2/templates"); err != nil {
		log.Fatalf("%v", err)
	}
}

type KubeObject map[string]interface{}
type KubeObjectList []KubeObject

func (o KubeObject) DeleteNamespace() {
	meta, ok := o["metadata"].(map[string]string)
	if !ok {
		panic("metadata bad")
	}
	delete(meta, "namespace")
	o["metadata"] = meta
}
func (o KubeObject) Filename() string {
	meta, ok := o["metadata"].(map[string]interface{})
	if !ok {
		panic("metadata bad")
	}
	nameVal, ok := meta["name"]
	if !ok {
		panic("metadata bad - name")
	}
	name, ok := nameVal.(string)
	if !ok {
		panic("name not string")
	}
	kindVal, ok := o["kind"]
	if !ok {
		panic("kind bad")
	}
	kind, ok := kindVal.(string)
	if !ok {
		panic("kind not string")
	}
	return strings.Replace(strcase.ToSnake(name+kind), "-", "_", -1)+".yaml"
}

func writeTemplates(dir string) error {
	objs, err := ParseKubeManifest(source)
	if err != nil {
		return err
	}
	os.MkdirAll(dir, 0755)
	for _, obj := range objs {
		yml, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(dir, obj.Filename()), yml, 0644)
	}
	return nil
}

func ParseKubeManifest(manifest string) (KubeObjectList, error) {
	snippets := strings.Split(manifest, "---")
	var objs KubeObjectList
	for _, objectYaml := range snippets {
		var kubeObj KubeObject
		if err := yaml.Unmarshal([]byte(objectYaml), &kubeObj); err != nil {
			return nil, err
		}
		if kubeObj == nil {
			log.Warnf("bad yaml: %v", objectYaml)
			continue
		}
		objs = append(objs, kubeObj)
	}
	return objs, nil
}

const source = `### Namespace ###
kind: Namespace
apiVersion: v1
metadata:
  name: linkerd

### Service Account Controller ###
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: linkerd-controller
  namespace: linkerd

### Controller RBAC ###
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: linkerd-linkerd-controller
rules:
- apiGroups: ["extensions", "apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["list", "get", "watch"]
- apiGroups: [""]
  resources: ["pods", "endpoints", "services", "namespaces", "replicationcontrollers"]
  verbs: ["list", "get", "watch"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: linkerd-linkerd-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: linkerd-linkerd-controller
subjects:
- kind: ServiceAccount
  name: linkerd-controller
  namespace: linkerd

### Service Account Prometheus ###
---
kind: ServiceAccount
apiVersion: v1
metadata:
  name: linkerd-prometheus
  namespace: linkerd

### Prometheus RBAC ###
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: linkerd-linkerd-prometheus
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: linkerd-linkerd-prometheus
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: linkerd-linkerd-prometheus
subjects:
- kind: ServiceAccount
  name: linkerd-prometheus
  namespace: linkerd

### Controller ###
---
kind: Service
apiVersion: v1
metadata:
  name: api
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: controller
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
spec:
  type: ClusterIP
  selector:
    linkerd.io/control-plane-component: controller
  ports:
  - name: http
    port: 8085
    targetPort: 8085

---
kind: Service
apiVersion: v1
metadata:
  name: proxy-api
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: controller
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
spec:
  type: ClusterIP
  selector:
    linkerd.io/control-plane-component: controller
  ports:
  - name: grpc
    port: 8086
    targetPort: 8086

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
  creationTimestamp: null
  labels:
    linkerd.io/control-plane-component: controller
  name: controller
  namespace: linkerd
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        linkerd.io/created-by: linkerd/cli stable-2.0.0
        linkerd.io/proxy-version: stable-2.0.0
      creationTimestamp: null
      labels:
        linkerd.io/control-plane-component: controller
        linkerd.io/control-plane-ns: linkerd
        linkerd.io/proxy-deployment: controller
    spec:
      containers:
      - args:
        - public-api
        - -prometheus-url=http://prometheus.linkerd.svc.cluster.local:9090
        - -controller-namespace=linkerd
        - -log-level=info
        image: gcr.io/linkerd-io/controller:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /ping
            port: 9995
          initialDelaySeconds: 10
        name: public-api
        ports:
        - containerPort: 8085
          name: http
        - containerPort: 9995
          name: admin-http
        readinessProbe:
          failureThreshold: 7
          httpGet:
            path: /ready
            port: 9995
        resources: {}
      - args:
        - destination
        - -enable-tls=false
        - -log-level=info
        image: gcr.io/linkerd-io/controller:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /ping
            port: 9999
          initialDelaySeconds: 10
        name: destination
        ports:
        - containerPort: 8089
          name: grpc
        - containerPort: 9999
          name: admin-http
        readinessProbe:
          failureThreshold: 7
          httpGet:
            path: /ready
            port: 9999
        resources: {}
      - args:
        - proxy-api
        - -addr=:8086
        - -log-level=info
        image: gcr.io/linkerd-io/controller:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /ping
            port: 9996
          initialDelaySeconds: 10
        name: proxy-api
        ports:
        - containerPort: 8086
          name: grpc
        - containerPort: 9996
          name: admin-http
        readinessProbe:
          failureThreshold: 7
          httpGet:
            path: /ready
            port: 9996
        resources: {}
      - args:
        - tap
        - -log-level=info
        - -controller-namespace=linkerd
        image: gcr.io/linkerd-io/controller:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /ping
            port: 9998
          initialDelaySeconds: 10
        name: tap
        ports:
        - containerPort: 8088
          name: grpc
        - containerPort: 9998
          name: admin-http
        readinessProbe:
          failureThreshold: 7
          httpGet:
            path: /ready
            port: 9998
        resources: {}
      - env:
        - name: LINKERD2_PROXY_LOG
          value: warn,linkerd2_proxy=info
        - name: LINKERD2_PROXY_BIND_TIMEOUT
          value: 10s
        - name: LINKERD2_PROXY_CONTROL_URL
          value: tcp://localhost.:8086
        - name: LINKERD2_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: LINKERD2_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: LINKERD2_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: LINKERD2_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: LINKERD2_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/linkerd-io/proxy:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        name: linkerd-proxy
        ports:
        - containerPort: 4143
          name: linkerd-proxy
        - containerPort: 4191
          name: linkerd-metrics
        readinessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        resources: {}
        securityContext:
          runAsUser: 2102
        terminationMessagePolicy: FallbackToLogsOnError
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/linkerd-io/proxy-init:stable-2.0.0
        imagePullPolicy: IfNotPresent
        name: linkerd-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
        terminationMessagePolicy: FallbackToLogsOnError
      serviceAccount: linkerd-controller
status: {}
---
kind: Service
apiVersion: v1
metadata:
  name: web
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: web
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
spec:
  type: ClusterIP
  selector:
    linkerd.io/control-plane-component: web
  ports:
  - name: http
    port: 8084
    targetPort: 8084
  - name: admin-http
    port: 9994
    targetPort: 9994

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
  creationTimestamp: null
  labels:
    linkerd.io/control-plane-component: web
  name: web
  namespace: linkerd
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        linkerd.io/created-by: linkerd/cli stable-2.0.0
        linkerd.io/proxy-version: stable-2.0.0
      creationTimestamp: null
      labels:
        linkerd.io/control-plane-component: web
        linkerd.io/control-plane-ns: linkerd
        linkerd.io/proxy-deployment: web
    spec:
      containers:
      - args:
        - -api-addr=api.linkerd.svc.cluster.local:8085
        - -static-dir=/dist
        - -template-dir=/templates
        - -uuid=e77de2f9-c6e0-4760-9d18-ec4d9e716ae6
        - -controller-namespace=linkerd
        - -log-level=info
        image: gcr.io/linkerd-io/web:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /ping
            port: 9994
          initialDelaySeconds: 10
        name: web
        ports:
        - containerPort: 8084
          name: http
        - containerPort: 9994
          name: admin-http
        readinessProbe:
          failureThreshold: 7
          httpGet:
            path: /ready
            port: 9994
        resources: {}
      - env:
        - name: LINKERD2_PROXY_LOG
          value: warn,linkerd2_proxy=info
        - name: LINKERD2_PROXY_BIND_TIMEOUT
          value: 10s
        - name: LINKERD2_PROXY_CONTROL_URL
          value: tcp://proxy-api.linkerd.svc.cluster.local:8086
        - name: LINKERD2_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: LINKERD2_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: LINKERD2_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: LINKERD2_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: LINKERD2_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/linkerd-io/proxy:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        name: linkerd-proxy
        ports:
        - containerPort: 4143
          name: linkerd-proxy
        - containerPort: 4191
          name: linkerd-metrics
        readinessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        resources: {}
        securityContext:
          runAsUser: 2102
        terminationMessagePolicy: FallbackToLogsOnError
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/linkerd-io/proxy-init:stable-2.0.0
        imagePullPolicy: IfNotPresent
        name: linkerd-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
        terminationMessagePolicy: FallbackToLogsOnError
status: {}
---
kind: Service
apiVersion: v1
metadata:
  name: prometheus
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: prometheus
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
spec:
  type: ClusterIP
  selector:
    linkerd.io/control-plane-component: prometheus
  ports:
  - name: admin-http
    port: 9090
    targetPort: 9090

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
  creationTimestamp: null
  labels:
    linkerd.io/control-plane-component: prometheus
  name: prometheus
  namespace: linkerd
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        linkerd.io/created-by: linkerd/cli stable-2.0.0
        linkerd.io/proxy-version: stable-2.0.0
      creationTimestamp: null
      labels:
        linkerd.io/control-plane-component: prometheus
        linkerd.io/control-plane-ns: linkerd
        linkerd.io/proxy-deployment: prometheus
    spec:
      containers:
      - args:
        - --storage.tsdb.retention=6h
        - --config.file=/etc/prometheus/prometheus.yml
        image: prom/prometheus:v2.4.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /-/healthy
            port: 9090
          initialDelaySeconds: 30
          timeoutSeconds: 30
        name: prometheus
        ports:
        - containerPort: 9090
          name: admin-http
        readinessProbe:
          httpGet:
            path: /-/ready
            port: 9090
          initialDelaySeconds: 30
          timeoutSeconds: 30
        resources: {}
        volumeMounts:
        - mountPath: /etc/prometheus
          name: prometheus-config
          readOnly: true
      - env:
        - name: LINKERD2_PROXY_LOG
          value: warn,linkerd2_proxy=info
        - name: LINKERD2_PROXY_BIND_TIMEOUT
          value: 10s
        - name: LINKERD2_PROXY_CONTROL_URL
          value: tcp://proxy-api.linkerd.svc.cluster.local:8086
        - name: LINKERD2_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: LINKERD2_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: LINKERD2_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: LINKERD2_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: LINKERD2_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: LINKERD2_PROXY_OUTBOUND_ROUTER_CAPACITY
          value: "10000"
        image: gcr.io/linkerd-io/proxy:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        name: linkerd-proxy
        ports:
        - containerPort: 4143
          name: linkerd-proxy
        - containerPort: 4191
          name: linkerd-metrics
        readinessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        resources: {}
        securityContext:
          runAsUser: 2102
        terminationMessagePolicy: FallbackToLogsOnError
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/linkerd-io/proxy-init:stable-2.0.0
        imagePullPolicy: IfNotPresent
        name: linkerd-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
        terminationMessagePolicy: FallbackToLogsOnError
      serviceAccount: linkerd-prometheus
      volumes:
      - configMap:
          name: prometheus-config
        name: prometheus-config
status: {}
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: prometheus-config
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: prometheus
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
data:
  prometheus.yml: |-
    global:
      scrape_interval: 10s
      scrape_timeout: 10s
      evaluation_interval: 10s

    scrape_configs:
    - job_name: 'prometheus'
      static_configs:
      - targets: ['localhost:9090']

    - job_name: 'grafana'
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names: ['linkerd']
      relabel_configs:
      - source_labels:
        - __meta_kubernetes_pod_container_name
        action: keep
        regex: ^grafana$

    - job_name: 'linkerd-controller'
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names: ['linkerd']
      relabel_configs:
      - source_labels:
        - __meta_kubernetes_pod_label_linkerd_io_control_plane_component
        - __meta_kubernetes_pod_container_port_name
        action: keep
        regex: (.*);admin-http$
      - source_labels: [__meta_kubernetes_pod_container_name]
        action: replace
        target_label: component

    - job_name: 'linkerd-proxy'
      kubernetes_sd_configs:
      - role: pod
      relabel_configs:
      - source_labels:
        - __meta_kubernetes_pod_container_name
        - __meta_kubernetes_pod_container_port_name
        - __meta_kubernetes_pod_label_linkerd_io_control_plane_ns
        action: keep
        regex: ^linkerd-proxy;linkerd-metrics;linkerd$
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: pod
      # special case k8s' "job" label, to not interfere with prometheus' "job"
      # label
      # __meta_kubernetes_pod_label_linkerd_io_proxy_job=foo =>
      # k8s_job=foo
      - source_labels: [__meta_kubernetes_pod_label_linkerd_io_proxy_job]
        action: replace
        target_label: k8s_job
      # __meta_kubernetes_pod_label_linkerd_io_proxy_deployment=foo =>
      # deployment=foo
      - action: labelmap
        regex: __meta_kubernetes_pod_label_linkerd_io_proxy_(.+)
      # drop all labels that we just made copies of in the previous labelmap
      - action: labeldrop
        regex: __meta_kubernetes_pod_label_linkerd_io_proxy_(.+)
      # __meta_kubernetes_pod_label_linkerd_io_foo=bar =>
      # foo=bar
      - action: labelmap
        regex: __meta_kubernetes_pod_label_linkerd_io_(.+)

### Grafana ###
---
kind: Service
apiVersion: v1
metadata:
  name: grafana
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: grafana
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
spec:
  type: ClusterIP
  selector:
    linkerd.io/control-plane-component: grafana
  ports:
  - name: http
    port: 3000
    targetPort: 3000

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
  creationTimestamp: null
  labels:
    linkerd.io/control-plane-component: grafana
  name: grafana
  namespace: linkerd
spec:
  replicas: 1
  strategy: {}
  template:
    metadata:
      annotations:
        linkerd.io/created-by: linkerd/cli stable-2.0.0
        linkerd.io/proxy-version: stable-2.0.0
      creationTimestamp: null
      labels:
        linkerd.io/control-plane-component: grafana
        linkerd.io/control-plane-ns: linkerd
        linkerd.io/proxy-deployment: grafana
    spec:
      containers:
      - image: gcr.io/linkerd-io/grafana:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /api/health
            port: 3000
        name: grafana
        ports:
        - containerPort: 3000
          name: http
        readinessProbe:
          failureThreshold: 10
          httpGet:
            path: /api/health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 30
        resources: {}
        volumeMounts:
        - mountPath: /etc/grafana
          name: grafana-config
          readOnly: true
      - env:
        - name: LINKERD2_PROXY_LOG
          value: warn,linkerd2_proxy=info
        - name: LINKERD2_PROXY_BIND_TIMEOUT
          value: 10s
        - name: LINKERD2_PROXY_CONTROL_URL
          value: tcp://proxy-api.linkerd.svc.cluster.local:8086
        - name: LINKERD2_PROXY_CONTROL_LISTENER
          value: tcp://0.0.0.0:4190
        - name: LINKERD2_PROXY_METRICS_LISTENER
          value: tcp://0.0.0.0:4191
        - name: LINKERD2_PROXY_PRIVATE_LISTENER
          value: tcp://127.0.0.1:4140
        - name: LINKERD2_PROXY_PUBLIC_LISTENER
          value: tcp://0.0.0.0:4143
        - name: LINKERD2_PROXY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        image: gcr.io/linkerd-io/proxy:stable-2.0.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        name: linkerd-proxy
        ports:
        - containerPort: 4143
          name: linkerd-proxy
        - containerPort: 4191
          name: linkerd-metrics
        readinessProbe:
          httpGet:
            path: /metrics
            port: 4191
          initialDelaySeconds: 10
        resources: {}
        securityContext:
          runAsUser: 2102
        terminationMessagePolicy: FallbackToLogsOnError
      initContainers:
      - args:
        - --incoming-proxy-port
        - "4143"
        - --outgoing-proxy-port
        - "4140"
        - --proxy-uid
        - "2102"
        - --inbound-ports-to-ignore
        - 4190,4191
        image: gcr.io/linkerd-io/proxy-init:stable-2.0.0
        imagePullPolicy: IfNotPresent
        name: linkerd-init
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
          privileged: false
        terminationMessagePolicy: FallbackToLogsOnError
      volumes:
      - configMap:
          items:
          - key: grafana.ini
            path: grafana.ini
          - key: datasources.yaml
            path: provisioning/datasources/datasources.yaml
          - key: dashboards.yaml
            path: provisioning/dashboards/dashboards.yaml
          name: grafana-config
        name: grafana-config
status: {}
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: grafana-config
  namespace: linkerd
  labels:
    linkerd.io/control-plane-component: grafana
  annotations:
    linkerd.io/created-by: linkerd/cli stable-2.0.0
data:
  grafana.ini: |-
    instance_name = linkerd-grafana

    [server]
    root_url = %(protocol)s://%(domain)s:/api/v1/namespaces/linkerd/services/grafana:http/proxy/

    [auth]
    disable_login_form = true

    [auth.anonymous]
    enabled = true
    org_role = Editor

    [auth.basic]
    enabled = false

    [analytics]
    check_for_updates = false

  datasources.yaml: |-
    apiVersion: 1
    datasources:
    - name: prometheus
      type: prometheus
      access: proxy
      orgId: 1
      url: http://prometheus.linkerd.svc.cluster.local:9090
      isDefault: true
      jsonData:
        timeInterval: "5s"
      version: 1
      editable: true

  dashboards.yaml: |-
    apiVersion: 1
    providers:
    - name: 'default'
      orgId: 1
      folder: ''
      type: file
      disableDeletion: true
      editable: true
      options:
        path: /var/lib/grafana/dashboards
        homeDashboardId: linkerd-top-line
---
`
