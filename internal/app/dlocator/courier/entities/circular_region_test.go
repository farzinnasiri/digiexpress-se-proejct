package entities

import (
	"strings"
	"testing"
)

func TestConvertToMetersIfRadiusIsKiloMetres(t *testing.T) {
	cr := CircularRegion{Radius: 1}
	cr = cr.convertToMetersIfRadiusIsInKiloMetres("Kilometers")

	if cr.Radius != 1000 {
		t.Errorf("incorrect converstion, got %f expected %d", cr.Radius, 1000)
	}

}

func TestValidate(t *testing.T) {

}

func TestValidateCenterIsInBoundShouldReturnErrIfLatIsOutBound(t *testing.T) {
	cr := CircularRegion{
		CenterLatitude:  100,
		CenterLongitude: 50,
	}

	err := cr.validateCenterIsInBound()
	if err == nil {
		t.Errorf("should return non empty error")
		t.FailNow()
	}

	if !strings.Contains(err.Error(), "latitude") {
		t.Errorf("Should retun error explicitly declreading latitude is invalid")
	}

}

func TestValidateCenterIsInBoundShouldReturnErrIfLonIsOutBound(t *testing.T) {
	cr := CircularRegion{
		CenterLatitude:  33,
		CenterLongitude: 100,
	}

	err := cr.validateCenterIsInBound()
	if err == nil {
		t.Errorf("should return non empty error")
		t.FailNow()
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}

}

func TestValidateRadiusShouldReturnErrIfRadiusIsZero(t *testing.T) {
	cr := CircularRegion{
		Radius: 0,
	}

	err := cr.validateRadius()
	if err == nil {
		t.Errorf("should return non empty error")
		t.FailNow()
	}

	if !strings.Contains(err.Error(), "zero") {
		t.Errorf("Should retun error explicitly declreading zero is invalid")
	}
}

func TestValidateRadiusShouldReturnErrIfRadiusIsTooLarge(t *testing.T) {
	cr := CircularRegion{
		Radius: 2000 * 1000,
	}

	err := cr.validateRadius()
	if err == nil {
		t.Errorf("should return non empty error")
		t.FailNow()
	}

	if !strings.Contains(err.Error(), "greater") {
		t.Errorf("Should retun error explicitly declreading value is invalid")
	}
}

func TestValidateRadiusTypeShouldReturnErrIfTypeIsInvalid(t *testing.T) {
	cr := CircularRegion{}

	err := cr.validateRadiusType("miles")
	if err == nil {
		t.Errorf("should return non empty error")
	}

}

func TestValidateRadiusTypeShouldNotReturnErrIfTypeIsValid(t *testing.T) {
	cr := CircularRegion{}

	err := cr.validateRadiusType("Meters")
	if err != nil {
		t.Errorf("should  not return error: %s", err.Error())
	}

	err = cr.validateRadiusType("Kilometers")
	if err != nil {
		t.Errorf("should  not return error: %s", err.Error())
	}

}
