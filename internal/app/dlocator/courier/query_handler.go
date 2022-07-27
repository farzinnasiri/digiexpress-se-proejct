package courier

import (
	"context"

	"github.com/digiexpress/dlocator/internal/app/dlocator/courier/entities"
)

type QueryHandler struct {
	cacheRepository CacheRepository
}

func NewCourierQueryHandler(cacheRepository CacheRepository) QueryHandler {
	return QueryHandler{cacheRepository}
}

func (q QueryHandler) GetCouriersInsideCircle(ctx context.Context,
	query GetCouriersInsideCircleQuery) ([]Courier, error) {
	circularRegion, err := entities.NewCircularRegion(query.CenterLatitude, query.CenterLongitude,
		query.Radius, query.RadiusType)
	if err != nil {
		return nil, err
	}

	couriers, err := q.cacheRepository.GetCouriersInsideCircularRegion(CircularRegionToDTO(circularRegion))
	if err != nil {
		return nil, err
	}

	return couriers, nil
}
