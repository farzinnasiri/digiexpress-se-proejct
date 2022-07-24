package benchmarks

import (
	"math/rand"

	dv1 "github.com/digiexpress/dlocator/pkg/api/v1"
)

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomLocationInTehran() *dv1.Location {
	return &dv1.Location{
		Latitude:  randomFloat64(35.609668, 35.758597),
		Longitude: randomFloat64(51.236096, 51.492949),
	}
}
