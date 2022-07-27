package courier

import (
	"github.com/digiexpress/dlocator/internal/pkg/service"
)

type CacheRepositoryImpl struct {
	cacheService service.CacheService
}

const CouriersSet = "couriers"

func NewCacheRepository(cacheService service.CacheService) CacheRepositoryImpl {
	return CacheRepositoryImpl{cacheService}
}

func (c CacheRepositoryImpl) SaveCourierLocation(courier Courier) error {
	return c.cacheService.AddGeoLocation(
		CouriersSet, service.CourierDTO{
			DriverId:  courier.DriverID,
			Latitude:  courier.Location.Latitude,
			Longitude: courier.Location.Longitude,
		})
}

func (c CacheRepositoryImpl) GetCouriersInsideCircularRegion(circularRegionDTO CircularRegionDTO) ([]Courier, error) {
	searchByRadiusDTO := service.SearchByRadiusDTO{
		CenterLatitude:  circularRegionDTO.CenterLatitude,
		CenterLongitude: circularRegionDTO.CenterLongitude,
		Radius:          circularRegionDTO.Radius,
	}

	courierDTOs, err := c.cacheService.GeoSearchByRadius(CouriersSet, searchByRadiusDTO)
	if err != nil {
		return nil, err
	}

	couriers := make([]Courier, len(courierDTOs))
	for i, location := range courierDTOs {
		couriers[i] = Courier{
			DriverID: location.DriverId,
			Location: LocationDTO{
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
		}
	}

	return couriers, nil
}
