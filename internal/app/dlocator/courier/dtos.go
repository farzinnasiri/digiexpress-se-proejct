package courier

import (
	"github.com/digiexpress/dlocator/internal/app/dlocator/courier/entities"
)

type Courier struct {
	DriverID string
	Location LocationDTO
}

type LocationDTO struct {
	Latitude  float64
	Longitude float64
}

func LocationToDTO(location entities.Location) LocationDTO {
	return LocationDTO{
		Latitude:  location.GetLatitude(),
		Longitude: location.GetLongitude(),
	}
}

type CircularRegionDTO struct {
	CenterLatitude  float64
	CenterLongitude float64
	Radius          float64
}

func CircularRegionToDTO(circularRegion entities.CircularRegion) CircularRegionDTO {
	return CircularRegionDTO{
		CenterLatitude:  circularRegion.CenterLatitude,
		CenterLongitude: circularRegion.CenterLongitude,
		Radius:          circularRegion.Radius,
	}
}
