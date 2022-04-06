module github.com/gojekfarm/envoy-lb-operator

go 1.12

require (
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/gogo/protobuf v1.2.1
	github.com/gojektech/kubehandler v0.0.0-20190321033534-0bd438f7bfb9
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.36.0
	k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v10.0.0+incompatible
)
