package dlocator

import (
	"github.com/digiexpress/dlocator/internal/app/dlocator/courier"
	dv1 "github.com/digiexpress/dlocator/internal/pkg/api"
)

func InjectCourierLocatorServer(commandHandler courier.CommandHandler, queryHandler courier.QueryHandler) dv1.CourierLocatorServer {
	return NewCourierLocatorHandler(commandHandler, queryHandler)
}
