package idmapper

type result struct {
	value string
	found bool
}

type request struct {
	key      string
	response chan<- result
}

// LockFree object for mapping values' ID and name. May be shared between goroutines. Lock-free (mutex-free) version of IDMapper.
type LockFree struct {
	source   SourceReader
	reload   chan ValuesMap
	requests chan request
}

// NewLockFree creates new LockFree and load values using SourceReader
func NewLockFree(source SourceReader, done <-chan struct{}) (*LockFree, error) {
	idMapper := &LockFree{
		source:   source,
		reload:   make(chan ValuesMap),
		requests: make(chan request),
	}

	go idMapper.server(done)
	return idMapper, idMapper.Reload()
}

// Get gets value's name for given ID. Return value is pair of value's name and boolean if value was found
func (idMapper *LockFree) Get(id string) (string, bool) {
	response := make(chan result)
	idMapper.requests <- request{id, response}
	res := <-response
	return res.value, res.found
}

func (idMapper *LockFree) server(done <-chan struct{}) {
	values := make(ValuesMap)

	for {
		select {
		case req := <-idMapper.requests:
			value, found := values[req.key]
			req.response <- result{
				value: value,
				found: found,
			}
		case reloadedValues := <-idMapper.reload:
			values = reloadedValues
		case <-done:
			return
		}
	}
}

// Reload reloads id mapper values using SourceReader
func (idMapper *LockFree) Reload() error {
	newValues, err := idMapper.source.Read()
	if err != nil {
		return err
	}

	idMapper.reload <- newValues
	return nil
}
