package envoy_test

import (
	"testing"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/stretchr/testify/assert"
)

func TestNilNodeID(t *testing.T) {
	assert.Equal(t, "unknown", envoy.Hasher{}.ID(nil))
}

func TestHasherNodeID(t *testing.T) {
	id := "12736478236847623874628"
	assert.Equal(t, id, envoy.Hasher{}.ID(&core.Node{Id: id}))
}
