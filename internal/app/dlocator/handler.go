package dlocator

import (
	"context"

	"github.com/digiexpress/dlocator/internal/app/dlocator/courier"
	"github.com/digiexpress/dlocator/internal/pkg/api"
)

type CourierLocatorHandler struct {
	api.CourierLocatorServer

	courier.CommandHandler
	courier.QueryHandler
}

func NewCourierLocatorHandler(commandHandler courier.CommandHandler,
	queryHandler courier.QueryHandler) *CourierLocatorHandler {
	return &CourierLocatorHandler{CommandHandler: commandHandler, QueryHandler: queryHandler}
}

func (c *CourierLocatorHandler) PushGpsPoints(stream api.CourierLocator_PushGpsPointsServer) error {
	for {
		courierGpsPoint, err := stream.Recv()
		if err != nil {
			return err
		}
		streamContext := stream.Context()

		pushCourierLocationCommand := courier.PushCourierLocation{
			DriverID: courierGpsPoint.GetDriverId(),
			Location: courier.LocationDTO{
				Latitude:  courierGpsPoint.GetLocation().GetLatitude(),
				Longitude: courierGpsPoint.GetLocation().GetLongitude(),
			},
		}

		// todo error handling
		if err = c.CommandHandler.PushCourierLocation(streamContext, pushCourierLocationCommand); err != nil {
			return err
		}
	}
}

func (c *CourierLocatorHandler) FindNearbyCouriers(ctx context.Context,
	request *api.FindNearbyCouriersQuery,
) (*api.CourierList, error) {
	query := courier.GetCouriersInsideCircleQuery{
		CenterLatitude:  request.GetTarget().GetLatitude(),
		CenterLongitude: request.GetTarget().GetLongitude(),
		Radius:          request.GetRadius().Amount,
		RadiusType:      request.GetRadius().GetUnit().String(),
	}

	couriers, err := c.QueryHandler.GetCouriersInsideCircle(ctx, query)
	if err != nil {
		return nil, err
	}

	response := &api.CourierList{}
	points := make([]*api.CourierGpsPoint, 0)
	for _, co := range couriers {
		points = append(points, &api.CourierGpsPoint{
			DriverId: co.DriverID,
			Location: &api.Location{
				Latitude:  co.Location.Latitude,
				Longitude: co.Location.Longitude,
			},
		})
	}

	response.Points = points

	return response, nil
}
