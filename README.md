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
