package courier

type GetCouriersInsideCircleQuery struct {
	CenterLatitude  float64
	CenterLongitude float64
	Radius          float64
	RadiusType      string
}
