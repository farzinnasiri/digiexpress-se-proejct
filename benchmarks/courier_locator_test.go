//go:build exclud

package benchmarks

import (
	"context"
	"fmt"
	dv1 "github.com/digiexpress/dlocator/internal/pkg/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"testing"
	"time"
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
	//
	cancel()
}

func TestFindNearByCouriersShouldReturnOnlyCouriersInsideRegion(t *testing.T) {
	rand.Seed(100)

	couriersInside := []*Courier{
		{
			id: "1",
			location: &dv1.Location{
				Latitude:  35.82871,
				Longitude: 50.95665,
			},
		},
		{
			id: "2",
			location: &dv1.Location{
				Latitude:  35.85070,
				Longitude: 50.95047,
			},
		},
		{
			id: "3",
			location: &dv1.Location{
				Latitude:  35.86795,
				Longitude: 50.97450,
			},
		},
		{
			id: "4",
			location: &dv1.Location{
				Latitude:  35.83622,
				Longitude: 50.91545,
			},
		},
		{
			id: "5",
			location: &dv1.Location{
				Latitude:  35.84457,
				Longitude: 51.00815,
			},
		},
	}

	couriersOutside := []*Courier{
		{
			id: "6",
			location: &dv1.Location{
				Latitude:  35.71132,
				Longitude: 50.93318,
			},
		},
		{
			id: "7",
			location: &dv1.Location{
				Latitude:  35.79193,
				Longitude: 50.82420,
			},
		},
		{
			id: "8",
			location: &dv1.Location{
				Latitude:  35.88214,
				Longitude: 51.06607,
			},
		},
	}

	couriers := append(couriersInside, couriersOutside...)

	ctx, cancel := context.WithCancel(context.Background())
	server := setupServer(ctx)
	client := setupClient(ctx, server)

	simulator := &CourierMovementSimulation{
		driverGpsPushRate: time.Second * 1,
		couriers:          couriers,
	}

	go simulator.RunCouriersWithConstantLocation(ctx, client)
	time.Sleep(time.Second * 1)

	nearbyCouriersList, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
		Target: &dv1.Location{
			Latitude:  35.84334,
			Longitude: 50.96795,
		},
		Radius: &dv1.Distance{Amount: 8, Unit: dv1.Unit_Kilometers},
	})
	if err != nil {
		t.FailNow()
	}

	if len(nearbyCouriersList.Points) != 5 {
		t.Errorf("list should contain 5 couriers")
	}

	for _, courier := range couriersInside {
		if !(ContainsId(nearbyCouriersList.Points, courier.id)) {
			t.Errorf(fmt.Sprintf("list should contain courier with id %s", courier.id))
		}
	}

	cancel()
}

// This test needs a large enough TTL to work (TTL >= 2s)
func TestFindNearByCouriersShouldReturnSameListAfterAppRestart(t *testing.T) {
	couriersInside := []*Courier{
		{
			id: "1",
			location: &dv1.Location{
				Latitude:  35.82871,
				Longitude: 50.95665,
			},
		},
		{
			id: "2",
			location: &dv1.Location{
				Latitude:  35.85070,
				Longitude: 50.95047,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := setupServer(ctx)
	client := setupClient(ctx, server)

	simulator := &CourierMovementSimulation{
		driverGpsPushRate: time.Second * 1,
		couriers:          couriersInside,
	}

	go simulator.RunCouriersWithConstantLocation(ctx, client)
	time.Sleep(time.Second * 1)

	nearbyCouriersList, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
		Target: &dv1.Location{
			Latitude:  35.84334,
			Longitude: 50.96795,
		},
		Radius: &dv1.Distance{Amount: 8, Unit: dv1.Unit_Kilometers},
	})
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	cancel()

	ctx, cancel = context.WithCancel(context.Background())
	server = setupServer(ctx)
	client = setupClient(ctx, server)

	nearbyCouriersListAfterRestart, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
		Target: &dv1.Location{
			Latitude:  35.84334,
			Longitude: 50.96795,
		},
		Radius: &dv1.Distance{Amount: 8, Unit: dv1.Unit_Kilometers},
	})
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	if !isEqual(nearbyCouriersList, nearbyCouriersListAfterRestart) {
		t.Errorf("couriers should not change after restart")
	}

	cancel()
}

// It's better to manually set a short ttl for this test to run fast
func TestFindNearByCouriersShouldNotReturnAnyCourierAfterTTLisExpired(t *testing.T) {
	couriersInside := []*Courier{
		{
			id: "1",
			location: &dv1.Location{
				Latitude:  35.82871,
				Longitude: 50.95665,
			},
		},
		{
			id: "2",
			location: &dv1.Location{
				Latitude:  35.85070,
				Longitude: 50.95047,
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := setupServer(ctx)
	client := setupClient(ctx, server)

	simulator := &CourierMovementSimulation{
		couriers: couriersInside,
	}

	go simulator.RunCouriersOnlyOnce(ctx, client)

	ttl, err := readTTLFromConfig()
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	time.Sleep(time.Second * time.Duration(ttl))

	nearbyCouriersListAfterRestart, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
		Target: &dv1.Location{
			Latitude:  35.84334,
			Longitude: 50.96795,
		},
		Radius: &dv1.Distance{Amount: 8, Unit: dv1.Unit_Kilometers},
	})
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}

	if len(nearbyCouriersListAfterRestart.Points) != 0 {
		t.Errorf(fmt.Sprintf("no courier should be returend, got %v couriers",
			len(nearbyCouriersListAfterRestart.Points)))
	}

	cancel()
}

// use this function for testing already deployed server
func BenchmarkServer(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("localhost:%d", 9000), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panicf("could not dial grpc server on port %d", 9000)
	}

	client := dv1.NewCourierLocatorClient(conn)
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
	//
	cancel()
}

type AppConfig struct {
	Redis RedisConfig
}

type RedisConfig struct {
	TTL int // in seconds
}

func readTTLFromConfig() (int, error) {
	var config AppConfig
	viper.SetConfigFile("./configs/config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return 0, fmt.Errorf("faild to read file config: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return 0, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config.Redis.TTL, nil
}

func isEqual(list1 *dv1.CourierList, list2 *dv1.CourierList) bool {
	points1 := list1.Points
	points2 := list2.Points

	for i, point1 := range points1 {
		if point1.DriverId != points2[i].DriverId {
			return false
		} else if point1.Location.Latitude != points2[i].Location.Latitude {
			return false
		} else if point1.Location.Longitude != points2[i].Location.Longitude {
			return false
		}
	}
	return true
}

func ContainsId(s []*dv1.CourierGpsPoint, id string) bool {
	for _, v := range s {
		if v.DriverId == id {
			return true
		}
	}

	return false
}
