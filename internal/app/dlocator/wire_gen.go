// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package dlocator

import (
	"github.com/digiexpress/dlocator/internal/pkg/api"
	"github.com/digiexpress/dlocator/internal/pkg/config"
)

// Injectors from wire.go:

func CreateApp() (*App, error) {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		return nil, err
	}
	courierLocatorServer := InjectCourierLocatorServer()
	dLocatorGrpcServer, err := api.NewDLocatorGrpcServer(courierLocatorServer, appConfig)
	if err != nil {
		return nil, err
	}
	app := NewApp(appConfig, dLocatorGrpcServer)
	return app, nil
}
