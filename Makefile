all: clean	test	build
APP=envoy-lb-operator
ALL_PACKAGES=$(shell go list ./...)
SOURCE_DIRS=$(shell go list ./... | cut -d "/" -f4 | uniq)

clean:
	rm -rf ./out
	GO111MODULE=on go mod tidy -v

setup:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/fzipp/gocyclo

check-quality: lint fmt vet

build:
	@echo "Building './out/envoy-lb-operator'..."
	@mkdir -p ./out
	@go build -o out/envoy-lb-operator

test:
	go test -v ./...

fmt:
	gofmt -l -s -w $(SOURCE_DIRS)

imports:
	goimports -l -w -v $(SOURCE_DIRS)

cyclo:
	gocyclo -over 7 $(SOURCE_DIRS)

vet:
	go vet ./...

lint:
	@if [[ `golint $(ALL_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } | wc -l | tr -d ' '` -ne 0 ]]; then \
          golint $(ALL_PACKAGES) | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }; \
          exit 2; \
    fi;

linux-build:
	env GOOS=linux go build -o out/$(APP)-linux
docker-build: linux-build
	docker build --no-cache -t $(APP) . ;\

docker-dev: docker-build
	kubectl config use-context docker-for-desktop
	$(MAKE) kube-deploy

kube-del:
	helm delete --purge lb | true
	kubectl delete deployment/$(APP) | true
	kubectl delete svc/$(APP) | true

kube-deploy: kube-del
	kubectl apply -f dev/$(APP)-deployment.yaml
	kubectl apply -f dev/$(APP)-svc.yaml
	helm install --name lb stable/envoy -f ./dev/envoy-values.yaml

minikube-dev:
	@eval $$(SHELL=bash minikube docker-env) ;\
	$(MAKE) docker-build
	kubectl config use-context minikube
	$(MAKE) kube-deploy
