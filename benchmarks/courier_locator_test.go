package benchmarks

import (
	"context"
	"math/rand"
	"testing"

	dv1 "github.com/digiexpress/dlocator/pkg/api/v1"
)

func BenchmarkDLocatorApiWhenStreaming(b *testing.B) {
	rand.Seed(100)
	ctx, cancel := context.WithCancel(context.Background())
	server := setupServer(ctx)
	client := setupClient(ctx, server)
	sim := NewCourierMovementSimulation(20_000)
	go sim.Run(ctx, client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
			Target: randomLocationInTehran(),
			Radius: &dv1.Distance{Amount: randomFloat64(50, 500), Unit: dv1.Unit_Meters},
		})
		if err != nil {
			b.Fail()
		}
	}

	cancel()
}
