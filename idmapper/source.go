package idmapper

// SourceReader interface for reading values from source
type SourceReader interface {
	Read() (ValuesMap, error)
}

// SourceReaderFunc is adapter to allow use ordinary function as SourceReader
type SourceReaderFunc func() (ValuesMap, error)

func (fn SourceReaderFunc) Read() (ValuesMap, error) {
	return fn()
}
