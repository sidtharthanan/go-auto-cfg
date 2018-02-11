package main

import (
	"testing"

	t1 "github.com/sidtharthanan/go-auto-cfg/test_1"
	t2 "github.com/sidtharthanan/go-auto-cfg/test_2"
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

func TestOptionalConfig(t *testing.T) {
	t2Instance1 := t2.New()
	t2Instance1.Load("config2", "test_samples")

	assert.Equal(t, 25, t2Instance1.POOL_SIZE())
	assert.Equal(t, "localhost", t2Instance1.APP_HOST())
	assert.Equal(t, 0, t2Instance1.SIZE())
}
