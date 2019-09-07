package idmapper_test

import (
	"fmt"
	"testing"

	"github.com/danielkraic/idmapper/idmapper"
)

type result struct {
	key, value string
	found      bool
}

type benchmarkSource struct {
	ValuesCount int
}

func (bs benchmarkSource) Read() (idmapper.ValuesMap, error) {
	result := idmapper.ValuesMap{}

	if bs.ValuesCount < 1 {
		panic(fmt.Sprintf("benchmarkSource: valuesCount %d must be positive number", bs.ValuesCount))
	}
	for i := 0; i < bs.ValuesCount; i++ {
		key := fmt.Sprintf("%d", i)
		value := fmt.Sprintf("%d", i*2)
		result[key] = value
	}

	return result, nil
}

// Benchmark between IDMapper and LockFree
func BenchmarkIdMapperGet(b *testing.B) {
	for _, valuesCount := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%ditems", valuesCount), func(b *testing.B) {
			b.StopTimer()

			idMapper, err := idmapper.NewIDMapper(&benchmarkSource{ValuesCount: valuesCount})
			if err != nil {
				b.Fatal(err)
			}

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				idMapper.Get(fmt.Sprintf("%d", i))
			}
		})

		b.Run(fmt.Sprintf("LockFree_%ditems", valuesCount), func(b *testing.B) {
			b.StopTimer()

			done := make(chan struct{})
			idMapper, err := idmapper.NewLockFree(&benchmarkSource{ValuesCount: valuesCount}, done)
			if err != nil {
				b.Fatal(err)
			}

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				idMapper.Get(fmt.Sprintf("%d", i))
			}

			done <- struct{}{}
		})
	}
}
