package dlocator

import (
	"context"

	dv1 "github.com/digiexpress/dlocator/pkg/api/v1"
)

type CourierLocatorHandler struct{}

func NewCourierLocatorHandler() *CourierLocatorHandler {
	return &CourierLocatorHandler{}
}

var _ dv1.CourierLocatorServer = (*CourierLocatorHandler)(nil)

func (c *CourierLocatorHandler) PushGpsPoints(s dv1.CourierLocator_PushGpsPointsServer) error {
	panic("implement me")
}

func (c *CourierLocatorHandler) FindNearbyCouriers(_ context.Context, _ *dv1.FindNearbyCouriersQuery) (*dv1.CourierList, error) {
	panic("implement me")
}
