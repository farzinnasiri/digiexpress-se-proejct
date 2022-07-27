package entities

import (
	"fmt"

	"github.com/pkg/errors"
)

type CircularRegion struct {
	CenterLatitude  float64
	CenterLongitude float64
	Radius          float64
}

func NewCircularRegion(centerLatitude float64, centerLongitude float64,
	radius float64, radiusType string) (CircularRegion, error) {
	circularRegion := CircularRegion{}

	if err := circularRegion.validateRadiusType(radiusType); err != nil {
		return CircularRegion{}, err
	}

	circularRegion = circularRegion.fill(centerLatitude, centerLongitude, radius)
	circularRegion = circularRegion.convertToMetersIfRadiusIsInKiloMetres(radiusType)

	if err := circularRegion.validate(); err != nil {
		return CircularRegion{}, err
	}

	return circularRegion, nil
}

func (c CircularRegion) validateRadiusType(radiusType string) error {
	if radiusType != "Meters" && radiusType != "Kilometers" {
		return fmt.Errorf(
			"invalid radius type %s, expected %s or %s",
			radiusType, "Meters", "Kilometers")
	}

	return nil
}

func (c CircularRegion) convertToMetersIfRadiusIsInKiloMetres(radiusType string) CircularRegion {
	if radiusType == "Kilometers" {
		c.Radius *= 1_000
	}

	return c
}

func (c CircularRegion) fill(centerLatitude float64, centerLongitude float64, radius float64) CircularRegion {
	c.CenterLatitude = centerLatitude
	c.CenterLongitude = centerLongitude
	c.Radius = radius

	return c
}

func (c CircularRegion) validate() error {
	if err := c.validateCenterIsInBound(); err != nil {
		return err
	}

	if err := c.validateRadius(); err != nil {
		return err
	}
	return nil
}

func (c CircularRegion) validateCenterIsInBound() error {
	if c.CenterLatitude <= 25 || c.CenterLatitude >= 40 {
		return fmt.Errorf("latitude is out of IRAN's"+
			"boundery (25 to 40 deg), Location: {%f %f}", c.CenterLongitude, c.CenterLongitude)
	}

	if c.CenterLongitude <= 44 || c.CenterLongitude >= 63 {
		return fmt.Errorf("longitude is out of IRAN's"+
			"boundery (44 to 63 deg), Location: {%f %f}", c.CenterLatitude, c.CenterLongitude)
	}

	return nil
}

func (c CircularRegion) validateRadius() error {
	if c.Radius == 0 {
		return errors.New("Radius can't be zero")
	}

	if c.Radius > 1_000_000 {
		return errors.New(fmt.Sprintf("Radius can't be greater than 1000KM(1000000M), got %f", c.Radius))
	}
	return nil
}
