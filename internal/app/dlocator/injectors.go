package dlocator

import (
	dv1 "github.com/digiexpress/dlocator/pkg/api/v1"
)

func InjectCourierLocatorServer() dv1.CourierLocatorServer {
	return NewCourierLocatorHandler()
}
