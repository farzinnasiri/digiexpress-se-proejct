//go:build wireinject
// +build wireinject

package dlocator

import (
	"github.com/digiexpress/dlocator/internal/app/dlocator/courier"
	"github.com/digiexpress/dlocator/internal/pkg/api"
	"github.com/digiexpress/dlocator/internal/pkg/config"
	"github.com/digiexpress/dlocator/internal/pkg/service"
	"github.com/google/wire"
)

func CreateApp() (*App, error) {
	panic(
		wire.Build(
			config.NewAppConfig,
			api.NewDLocatorGrpcServer,
			NewApp,
			service.NewCacheService,
			wire.Bind(new(service.CacheService), new(service.CacheServiceImpl)),
			courier.NewCacheRepository,
			wire.Bind(new(courier.CacheRepository), new(courier.CacheRepositoryImpl)),
			courier.NewCourierQueryHandler,
			courier.NewCourierCommandHandler,
			InjectCourierLocatorServer,
		),
	)
}
