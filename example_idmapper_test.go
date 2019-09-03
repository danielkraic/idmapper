package idmapper_test

import (
	"fmt"
	"testing"

	"github.com/danielkraic/idmapper"
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

func Example(t *testing.T) {
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
