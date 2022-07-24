//go:build wireinject
// +build wireinject

package dlocator

import (
	"github.com/digiexpress/dlocator/internal/pkg/api"
	"github.com/digiexpress/dlocator/internal/pkg/config"
	"github.com/google/wire"
)

func CreateApp() (*App, error) {
	panic(
		wire.Build(
			config.NewAppConfig,
			api.NewDLocatorGrpcServer,
			NewApp,
			InjectCourierLocatorServer,
		),
	)
}
