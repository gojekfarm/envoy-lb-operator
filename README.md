# Envoy LB Operator  [![CircleCI](https://circleci.com/gh/gojekfarm/envoy-lb-operator.svg?style=svg)](https://circleci.com/gh/gojekfarm/envoy-lb-operator
)
Envoy LB Operator is an Envoy Control Plane for kubernetes.

# What
Configures Envoy as a Load balancer for a set of pods. This can be used in place of Kubernetes Service of Type LoadBalancer/Internal LoadBalancer and ClusterIP.
You get control on load balancing via fine grained control plane of Envoy.

# How
Envoy LB Operator discovers "Headless Services" with appropriate labels and injects them into envoy Control Plane.

### Headless Services
Headless Services are type `ClusterIp` with `clusterIP: None`
Here Kubernetes is hence used for discovery mechanism via dns.

### Envoy Discovery of Headless Services via envoy-lb-operator control plane.
1. Add Headless services with label `heritage: envoy-lb`
2. Additionally in case of grpc services `kubectl annotate svc/grpc-greeter envoy-lb-operator.gojektech.k8s.io/service-type=grpc`

Pods behind such services will get added directly to envoy as a `STRICT_DNS` cluster

# Why

1. This approach uses kubernetes for service discovery via headless services.
2. Uses Envoy for loadbalancer.
3. LoadBalancer is updated with new cluster by simply adding a new headless service. Operations on Envoy are minimized.


# Once
`kubectl create serviceaccount envoy-lb-operator`
`kubectl create clusterrolebinding envoy-lb-operator-binding --clusterrole=view --serviceaccount=default:envoy-lb-operator`

# Dev

Build
`make build`

Run Tests

`make test`

Build & Deploy on Monikube

`make minikube-dev`


Create deploys

`make kube-deploy`

Delete deploys

`make kube-del`

# Install 

This can be installed via a helm chart assuming you have `gojektech-incubator` repo. You can read about that here (https://github.com/gojektech/charts)

`helm install gojektech-incubator/envoy-lb-operator --name=some-envoy-cp`\

you can get the service details

`kubectl describe svc/some-envoy-cp-envoy-lb-operator`

This can then set the relevant values in `./examples/envoy.yaml`
Node ID should be appropriately set to match the config on Envoy and Operator
```
node:
  cluster: service_greeter
  id: id1 
```

This can support multiple envoy clusters. It can be configured as follows:

```
envoy_discovery_mapping:
  - envoy_id: "id1"
    upstream_label: "label1=val1"
    namespace: "namespace1"
  - envoy_id: "id2"
    upstream_label: "label2=val2"
    namespace: "namespace2"
```

All the upstreams having a given label (eg: label1=val1) will be registered to the corresponding envoy (eg: id1) in the specified namespace (eg: namespace1).


The Envoy should be pointed to the previously created service.
```
load_assignment:
  cluster_name: "xds_cluster"
  endpoints:
  - lb_endpoints:
      - endpoint:
          address:
            socket_address:
              address: some-cp-envoy-lb-operator
              port_value: 80
              protocol: "TCP"

```
Now Following will install Envoy pointing to the previously installed operator as Control plane.

`helm install --name some-lb  stable/envoy -f values.yaml`

Where `values.yaml` has the overridden  `files.envoy.yaml` value.

# Envoy connection refresh
Every time the snapshot version changes, envoy connections are dropped and recreated. The following configs determine the frequency of snapshot updates.

`refresh_interval_in_s` determines the time interval between envoy config updates.
`auto_refresh_conn` determines whether snaphot version should be incremented at an interval of `refresh_interval_in_s` irrespective of whether the routing config changes or not. If set to false, version will be incremented only when there are service updates.
```
envoy_discovery_mapping:
  - envoy_id: "id1"
    upstream_label: "label1=val1"
    endpoint_label: "label2=val2"
    namespace: "namespace1"
```
`endpoint_label` matches all the endpoints with the given label, and updates the snapshot when there are endpoint updates.