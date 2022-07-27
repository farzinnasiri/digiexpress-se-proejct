//go:build exclud

package benchmarks

import (
	"context"
	"time"

	dv1 "github.com/digiexpress/dlocator/internal/pkg/api"

	"github.com/google/uuid"
)

type Courier struct {
	id       string
	location *dv1.Location
}

func SpawnCourier() *Courier {
	return &Courier{
		id:       uuid.New().String(),
		location: randomLocationInTehran(),
	}
}

func (mc *Courier) Move(dLat float64, dLng float64) {
	mc.location.Latitude += randomFloat64(-dLat, dLat)
	mc.location.Longitude += randomFloat64(-dLng, dLng)
}

type CourierMovementSimulation struct {
	dLat              float64
	dLng              float64
	driverGpsPushRate time.Duration
	couriers          []*Courier
}

func NewCourierMovementSimulation(courierCount int) *CourierMovementSimulation {
	couriers := make([]*Courier, 0)
	for i := 0; i < courierCount; i++ {
		couriers = append(couriers, SpawnCourier())
	}

	return &CourierMovementSimulation{
		dLat:              0.0006,
		dLng:              0.0006,
		driverGpsPushRate: time.Second * 5,
		couriers:          couriers,
	}
}

func (s *CourierMovementSimulation) Run(ctx context.Context, client dv1.CourierLocatorClient) {
	stream, err := client.PushGpsPoints(ctx)
	if err != nil {
		panic("could not make a stream")
	}

	for {
		time.Sleep(s.driverGpsPushRate)
		for _, courier := range s.couriers {
			courier.Move(0.0006, 0.0006)
			err := stream.Send(&dv1.CourierGpsPoint{
				DriverId: courier.id,
				Location: courier.location,
			})
			if err != nil {
				return
			}
		}
	}
}

func (s *CourierMovementSimulation) RunCouriersWithConstantLocation(ctx context.Context, client dv1.CourierLocatorClient) {
	stream, err := client.PushGpsPoints(ctx)
	if err != nil {
		panic("could not make a stream")
	}

	for {
		for _, courier := range s.couriers {
			err = stream.Send(&dv1.CourierGpsPoint{
				DriverId: courier.id,
				Location: courier.location,
			})
			if err != nil {
				return
			}
		}
		time.Sleep(s.driverGpsPushRate)
	}
}

func (s *CourierMovementSimulation) RunCouriersOnlyOnce(ctx context.Context, client dv1.CourierLocatorClient) {
	stream, err := client.PushGpsPoints(ctx)
	if err != nil {
		panic("could not make a stream")
	}

	for {
		for _, courier := range s.couriers {
			err = stream.Send(&dv1.CourierGpsPoint{
				DriverId: courier.id,
				Location: courier.location,
			})
			if err != nil {
				return
			}
		}

		break
	}
}
