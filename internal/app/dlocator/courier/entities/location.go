package entities

import (
	"fmt"

	"github.com/pkg/errors"
)

type Location struct {
	latitude  float64
	longitude float64
}

func NewLocation(latitude float64, longitude float64) (Location, error) {
	location := Location{latitude: latitude, longitude: longitude}
	if err := location.validate(); err != nil {
		return Location{}, err
	}
	return location, nil
}

func (l Location) validate() error {
	if err := l.checkForZeroLocation(); err != nil {
		return err
	}
	if err := l.checkLocationIsInBound(); err != nil {
		return err
	}
	return nil
}

func (l Location) checkForZeroLocation() error {
	if l.latitude == 0 {
		return errors.New(fmt.Sprintf("latitude can't be 0, Location: {%f %f}", l.latitude, l.longitude))
	}

	if l.longitude == 0 {
		return errors.New(fmt.Sprintf("longitude can't be 0, Location: {%f %f}", l.latitude, l.longitude))
	}

	return nil
}

func (l Location) checkLocationIsInBound() error {
	if l.latitude <= 25 || l.latitude >= 40 {
		return errors.New(fmt.Sprintf("latitude is out of IRAN's"+
			"boundery (25 to 40 deg), Location: {%f %f}", l.latitude, l.longitude))
	}

	if l.longitude <= 44 || l.longitude >= 63 {
		return errors.New(fmt.Sprintf("longitude is out of IRAN's"+
			"boundery (44 to 63 deg), Location: {%f %f}", l.latitude, l.longitude))
	}

	return nil
}

func (l Location) GetLatitude() float64 {
	return l.latitude
}

func (l Location) GetLongitude() float64 {
	return l.longitude
}
