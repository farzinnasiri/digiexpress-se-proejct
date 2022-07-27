//go:build exclud

package benchmarks

import (
	"context"
	"fmt"
	"log"

	dv1 "github.com/digiexpress/dlocator/internal/pkg/api"

	"github.com/digiexpress/dlocator/internal/app/dlocator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupServer(ctx context.Context) *dlocator.App {
	appInstance, err := dlocator.CreateApp()
	if err != nil {
		log.Panicf("could not create dlocator application: %v", err)
	}
	go appInstance.Run(ctx)
	return appInstance
}

func setupClient(ctx context.Context, app *dlocator.App) dv1.CourierLocatorClient {
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("localhost:%d", app.Config.Grpc.Port), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panicf("could not dial grpc server on port %d", app.Config.Grpc.Port)
	}

	return dv1.NewCourierLocatorClient(conn)
}

//func TestServer(t *testing.T) dv1.CourierLocatorClient {
//	conn, err := grpc.DialContext(
//		ctx,
//		fmt.Sprintf("localhost:%d", 9000), grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	if err != nil {
//		log.Panicf("could not dial grpc server on port %d", 9000)
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	client := dv1.NewCourierLocatorClient(conn)
//	sim := NewCourierMovementSimulation(20_000)
//	go sim.Run(ctx, client)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, err := client.FindNearbyCouriers(ctx, &dv1.FindNearbyCouriersQuery{
//			Target: randomLocationInTehran(),
//			Radius: &dv1.Distance{Amount: randomFloat64(50, 500), Unit: dv1.Unit_Meters},
//		})
//		if err != nil {
//			b.Fail()
//		}
//	}
//	//
//	cancel()
//
//}
