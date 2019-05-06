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
