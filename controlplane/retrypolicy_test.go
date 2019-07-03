package controlplane_test

import (
	cp "github.com/gojekfarm/envoy-lb-operator/controlplane"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetryPolicy(t *testing.T) {
	retryOn := "xxx"
	retryPredicate := "test_retry_predicate"
	var numRetries uint32
	var hostSelectionMaxRetryAttempts int64
	numRetries = 10
	hostSelectionMaxRetryAttempts = 20

	retryPolicy := cp.RetryPolicy(retryOn, retryPredicate, numRetries, hostSelectionMaxRetryAttempts)

	assert.Equal(t, retryOn, retryPolicy.RetryOn)
	assert.Equal(t, retryPredicate, retryPolicy.RetryHostPredicate[0].Name)
	assert.Equal(t, numRetries, retryPolicy.NumRetries.Value)
	assert.Equal(t, hostSelectionMaxRetryAttempts, retryPolicy.HostSelectionRetryMaxAttempts)
}
