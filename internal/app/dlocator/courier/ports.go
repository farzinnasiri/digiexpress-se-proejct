package courier

type CacheRepository interface {
	SaveCourierLocation(courier Courier) error
	GetCouriersInsideCircularRegion(circularRegionDTO CircularRegionDTO) ([]Courier, error)
}
