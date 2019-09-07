package idmapper_test

import (
	"fmt"

	"github.com/danielkraic/idmapper/idmapper"
)

type TestingSource struct{}

func (ts *TestingSource) Read() (idmapper.ValuesMap, error) {
	return map[string]string{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}, nil
}

func SourceFunc() (idmapper.ValuesMap, error) {
	return map[string]string{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}, nil
}

func Example() {
	// idMapper, err := idmapper.NewIDMapper(idmapper.SourceReaderFunc(SourceFunc))
	idMapper, err := idmapper.NewIDMapper(&TestingSource{})
	if err != nil {
		panic(err)
	}

	printPair(idMapper, "a")
	printPair(idMapper, "A")
	printPair(idMapper, " c ")
}

func printPair(idMapper *idmapper.IDMapper, key string) {
	result, found := idMapper.Get(key)
	if found {
		fmt.Printf("id=%s, name=%s\n", key, result)
	} else {
		fmt.Printf("id=%s, NOT FOUND\n", key)
	}
}
