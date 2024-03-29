package idmapper_test

import (
	"fmt"
	"testing"

	"github.com/danielkraic/idmapper/idmapper"
	"github.com/stretchr/testify/assert"
)

var (
	errReadFailedString = "failed to read from source"
	errReadFailed       = fmt.Errorf(errReadFailedString)
)

type TestingSourceValid struct {
	values    idmapper.ValuesMap
	CallCount int
}

func (ts *TestingSourceValid) Read() (idmapper.ValuesMap, error) {
	ts.CallCount += 1
	return ts.values, nil
}

type TestingSourceInvalid struct{}

func (ts *TestingSourceInvalid) Read() (idmapper.ValuesMap, error) {
	return nil, errReadFailed
}

func TestNewIdMapperValid(t *testing.T) {
	source := &TestingSourceValid{
		values: idmapper.ValuesMap{},
	}
	_, err := idmapper.NewIDMapper(source)
	assert.Nil(t, err)
}

func TestNewIdMapperError(t *testing.T) {
	source := &TestingSourceInvalid{}
	_, err := idmapper.NewIDMapper(source)
	assert.EqualError(t, err, errReadFailedString)
}

func TestNewIdMapperValidFunc(t *testing.T) {

	fn := func() (idmapper.ValuesMap, error) {
		return idmapper.ValuesMap{
			"a": "A",
		}, nil
	}
	_, err := idmapper.NewIDMapper(idmapper.SourceReaderFunc(fn))
	assert.Nil(t, err)
}

func TestIdMapperGetExisting(t *testing.T) {
	values := idmapper.ValuesMap{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}

	source := &TestingSourceValid{
		values: values,
	}

	idMapper, err := idmapper.NewIDMapper(source)
	assert.Nil(t, err)

	for k, v := range values {
		result, found := idMapper.Get(k)
		assert.Equal(t, found, true)
		assert.Equal(t, result, v)
	}
}

func TestIdMapperGetNotExist(t *testing.T) {
	source := &TestingSourceValid{
		values: idmapper.ValuesMap{},
	}

	idMapper, err := idmapper.NewIDMapper(source)
	assert.Nil(t, err)

	nonExistingKeys := []string{"x", "y", "a ", " b"}
	for _, k := range nonExistingKeys {
		_, found := idMapper.Get(k)
		assert.Equal(t, found, false)
	}
}

func TestIdMapperReload(t *testing.T) {
	values := idmapper.ValuesMap{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}

	source := &TestingSourceValid{
		values: values,
	}

	idMapper, err := idmapper.NewIDMapper(source)
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
}
