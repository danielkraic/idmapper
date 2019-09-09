package idmapper_test

import (
	"testing"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/stretchr/testify/assert"
)

func TestNewLockFreeValid(t *testing.T) {
	source := &TestingSourceValid{
		values: idmapper.ValuesMap{},
	}

	done := make(chan struct{})
	_, err := idmapper.NewLockFree(source, done)
	assert.Nil(t, err)

	done <- struct{}{}
}

func TestNewLockFreeError(t *testing.T) {
	source := &TestingSourceInvalid{}

	done := make(chan struct{})
	_, err := idmapper.NewLockFree(source, done)
	assert.EqualError(t, err, errReadFailedString)

	done <- struct{}{}
}

func TestLockFreeGetExisting(t *testing.T) {
	values := idmapper.ValuesMap{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}

	source := &TestingSourceValid{
		values: values,
	}

	done := make(chan struct{})
	idMapper, err := idmapper.NewLockFree(source, done)
	assert.Nil(t, err)

	for k, v := range values {
		result, found := idMapper.Get(k)
		assert.Equal(t, found, true)
		assert.Equal(t, result, v)
	}

	done <- struct{}{}
}

func TestLockFreeGetNotExist(t *testing.T) {
	source := &TestingSourceValid{
		values: idmapper.ValuesMap{},
	}

	done := make(chan struct{})
	idMapper, err := idmapper.NewLockFree(source, done)
	assert.Nil(t, err)

	nonExistingKeys := []string{"x", "y", "a ", " b"}
	for _, k := range nonExistingKeys {
		_, found := idMapper.Get(k)
		assert.Equal(t, found, false)
	}

	done <- struct{}{}
}

func TestLockFreeReload(t *testing.T) {
	values := idmapper.ValuesMap{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}

	source := &TestingSourceValid{
		values: values,
	}

	done := make(chan struct{})
	idMapper, err := idmapper.NewLockFree(source, done)
	assert.Nil(t, err)

	for k, v := range values {
		result, found := idMapper.Get(k)
		assert.Equal(t, found, true)
		assert.Equal(t, result, v)
	}

	err = idMapper.Reload()
	assert.Nil(t, err)

	for k, v := range values {
		result, found := idMapper.Get(k)
		assert.Equal(t, found, true)
		assert.Equal(t, result, v)
	}

	assert.Equal(t, source.CallCount, 2)

	done <- struct{}{}
}
