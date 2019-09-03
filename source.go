package idmapper

// SourceReader interface for reading values from source
type SourceReader interface {
	Read() (ValuesMap, error)
}
