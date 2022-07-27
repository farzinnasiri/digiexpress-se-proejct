package service

type CacheService interface {
	GeoSearchByRadius(key string, queryByRadius SearchByRadiusDTO) ([]CourierDTO, error)
	AddGeoLocation(key string, courierDTO CourierDTO) error
}

type CourierDTO struct {
	DriverId           string
	Latitude           float64
	Longitude          float64
	DistanceFromCenter float64
}

type SearchByRadiusDTO struct {
	CenterLatitude  float64
	CenterLongitude float64
	Radius          float64
}
