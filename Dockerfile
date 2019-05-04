FROM golang:1.12-alpine
RUN apk update
RUN apk add git
WORKDIR /usr/src
ADD go.mod .
RUN go mod download
ADD $PWD /usr/src/envoy-lb-operator
WORKDIR /usr/src/envoy-lb-operator
RUN go build -o out/envoy-lb-operator

FROM alpine:3.9.2
WORKDIR /opt/envoy-lb-operator
COPY --from=0 usr/src/envoy-lb-operator/out/envoy-lb-operator .
RUN apk update
CMD ["./envoy-lb-operator", "server"]
