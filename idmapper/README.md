# idmapper

[![GoDoc](https://godoc.org/github.com/danielkraic/idmapper/idmapper?status.svg)](https://godoc.org/github.com/danielkraic/idmapper/idmapper)

IDMapper is in-memory cache for mapping IDs to Names. Can be used to cache lists of IDs and Names (eg country list, languages list).

## Documentation

https://godoc.org/github.com/danielkraic/idmapper/idmapper

## Example

```go
package main

import (
	"fmt"

	"github.com/danielkraic/idmapper/idmapper"
)

// structure implementing idmapper.SourceReader interface
type TestingSource struct{}

func (ts *TestingSource) Read() (idmapper.ValuesMap, error) {
    // can be function that reads data from DB or from http service

	return map[string]string{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}, nil
}

// Simple function can also be used as idmapper.SourceReader
func SourceFunc() (idmapper.ValuesMap, error) {
	return map[string]string{
		"":    "space",
		"a":   "A",
		"b":   "B",
		" c ": " C ",
	}, nil
}

func Example() {
    // create IDMapper
    idMapper, err := idmapper.NewIDMapper(&TestingSource{})
	if err != nil {
		panic(err)
	}

     // create IDMapper using function as SourceReader
    // idMapper, err := idmapper.NewIDMapper(idmapper.SourceReaderFunc(SourceFunc))
	
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
```

## IDMappers reloading

IDMappers can be reloaded automaticaly in background using [scheduler](https://github.com/danielkraic/idmapper/tree/master/scheduler)

```go
func idMappersReloaderExample() {
	loader := func() (idmapper.ValuesMap, error) {
		return idmapper.ValuesMap{
			"eur": "Euro",
			"usd": "US dollar",
			"czk": "Ceska koruna",
		}, nil
	}

	currency, err := idmapper.NewIDMapper(idmapper.SourceReaderFunc(loader))
	if err != nil {
		log.Fatal(err)
	}

	// create reloader
	reloader := scheduler.Scheduler{}
	// add job to reloader to reload currency IDMapper every second
	err = reloader.AddFunc(func() {
		failed := currency.Reload()
		if failed != nil {
			fmt.Printf("currency.Reload() failed: %s\n", failed)
		} else {
			fmt.Printf("currency.Reload() was successful\n")
		}
	}, 1*time.Second)

    // start reloader
	reloader.Start()
	defer reloader.Stop()

	for i := 0; i < 30; i++ {
		// get Name of 'eur' currency from IDMapper
		id := "eur"

		if name, found := currency.Get(id); found {
			fmt.Printf("id: %s, name: %s\n", id, name)
		} else {
			fmt.Printf("id: %s, NOT FOUND\n", id)
		}

		time.Sleep(100 * time.Millisecond)
	}

}
```