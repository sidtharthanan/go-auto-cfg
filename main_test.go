package main

import (
	"testing"

	t1 "github.com/sidtharthanan/go-auto-cfg/test_1"
	"github.com/stretchr/testify/assert"
)

func TestBasicDataTypes(t *testing.T) {
	t1.Load("config1", "test_samples")

	assert.Equal(t, 50, t1.POOL_SIZE())
	assert.Equal(t, []string{"staging", "production"}, t1.ENVS())
	assert.Equal(t, "busyqueue", t1.QUEUE_NAME())
	assert.Equal(t, true, t1.WORKER_ON())
	assert.Equal(t, 2.456, t1.SOME_CALC_FACTOR())
}

func TestMultipleInstances(t *testing.T) {
	tInstance1 := t1.New()
	tInstance1.Load("config1", "test_samples")

	tInstance2 := t1.New()
	tInstance2.Load("config1.2", "test_samples")

	assert.Equal(t, "busyqueue", tInstance1.QUEUE_NAME())
	assert.Equal(t, "freequeue", tInstance2.QUEUE_NAME())
}
