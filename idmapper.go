package idmapper

import "sync"

// ValuesMap map of values where map key is values' ID and map value is value's name
type ValuesMap map[string]string

// IDMapper object for mapping values IDs and names
type IDMapper struct {
	source SourceReader
	values ValuesMap
	mtx    sync.Mutex
}

// NewIDMapper creates new IDMapper and reload values using SourceReader
func NewIDMapper(source SourceReader) (*IDMapper, error) {
	idMapper := &IDMapper{
		source: source,
		values: make(ValuesMap),
	}

	return idMapper, idMapper.Reload()
}

// Get gets name for given ID. Return value is pair of value and boolean (if value was found)
func (idMapper *IDMapper) Get(id string) (string, bool) {
	idMapper.mtx.Lock()
	defer idMapper.mtx.Unlock()
	result, found := idMapper.values[id]
	return result, found
}

// Reload reloads id mapper values using SourceReader
func (idMapper *IDMapper) Reload() error {
	newValues, err := idMapper.source.Read()
	if err != nil {
		return err
	}

	idMapper.mtx.Lock()
	defer idMapper.mtx.Unlock()
	idMapper.values = newValues

	return nil
}
