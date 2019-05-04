package envoy_test

import (
	"testing"

	"github.com/gojekfarm/envoy-lb-operator/envoy"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotVersion(t *testing.T) {
	sn := envoy.NewSnapshot("node1")
	assert.Equal(t, int32(0), sn.Version)
}

func TestSnapshotVersionIncrementsOnStore(t *testing.T) {
	sn := envoy.NewSnapshot("node1")
	assert.Equal(t, int32(0), sn.Version)
	sn.Store()
	assert.Equal(t, int32(1), sn.Version)
	sn.Store()
	assert.Equal(t, int32(2), sn.Version)
}
