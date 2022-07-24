package benchmarks

import (
	"context"
	"fmt"
	"log"

	"github.com/digiexpress/dlocator/internal/app/dlocator"
	dv1 "github.com/digiexpress/dlocator/pkg/api/v1"
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
