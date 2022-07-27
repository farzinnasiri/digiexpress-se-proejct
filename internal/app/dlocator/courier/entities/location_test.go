package entities

import (
	"strings"
	"testing"
)

func TestCheckForZeroLocationShouldReturnErrIfLatIsZero(t *testing.T) {
	location := Location{
		latitude:  0,
		longitude: 35,
	}
	err := location.checkForZeroLocation()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "latitude") {
		t.Errorf("Should retun error explicitly declreading latitude is invalid")
	}
}

func TestCheckForZeroLocationShouldReturnErrIfLongIsZero(t *testing.T) {
	location := Location{
		latitude:  35,
		longitude: 0,
	}
	err := location.checkForZeroLocation()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}
}

func TestCheckLocationIsInBoundShouldReturnErrIfLatIsOutOfBound(t *testing.T) {
	// case 1
	location := Location{
		latitude:  100,
		longitude: 53,
	}
	err := location.checkLocationIsInBound()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "latitude") {
		t.Errorf("Should retun error explicitly declreading latitude is invalid")
	}

	// case 2
	location = Location{
		latitude:  10,
		longitude: 53,
	}
	err = location.checkLocationIsInBound()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "latitude") {
		t.Errorf("Should retun error explicitly declreading latitude is invalid")
	}
}

func TestCheckLocationIsInBoundShouldReturnErrIfLongIsOutOfBound(t *testing.T) {
	// case 1
	location := Location{
		latitude:  35,
		longitude: 10,
	}
	err := location.checkLocationIsInBound()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}

	// case 2
	location = Location{
		latitude:  35,
		longitude: 100,
	}
	err = location.checkLocationIsInBound()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}
}

func TestValidateShouldReturnErrIfLocationIsOutOfBound(t *testing.T) {
	location := Location{
		latitude:  35,
		longitude: 10,
	}
	err := location.validate()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}
}

func TestValidateShouldReturnErrIfLocationHasZeroValues(t *testing.T) {
	location := Location{
		latitude:  35,
		longitude: 0,
	}
	err := location.validate()
	if err == nil {
		t.Errorf("Should return non empty err")
	}

	if !strings.Contains(err.Error(), "longitude") {
		t.Errorf("Should retun error explicitly declreading longitude is invalid")
	}
}

func TestNewLocationWithValidInputs(t *testing.T) {
	_, err := NewLocation(35.0, 53.0)
	if err != nil {
		t.Errorf("Should not return error: %s", err)
	}
}

func TestNewLocationWithInvalidLat(t *testing.T) {
	_, err := NewLocation(0, 53.0)
	if err == nil {
		t.Errorf("Should not return error")
	}
}

func TestNewLocationWithInvalidLong(t *testing.T) {
	_, err := NewLocation(35.0, 0)
	if err == nil {
		t.Errorf("Should return error")
	}
}

func TestGetLatitudeShouldReturnLatIfInputsAreValid(t *testing.T) {
	testLat := 35.0
	location, _ := NewLocation(testLat, 53.0)
	if testLat != location.GetLatitude() {
		t.Errorf("latitudes don't match, expected %f got %f", testLat, location.GetLatitude())
	}
}

func TestGetLongitudeShouldReturnLongIfInputsAreValid(t *testing.T) {
	testLong := 53.0
	location, _ := NewLocation(35.0, testLong)
	if testLong != location.GetLongitude() {
		t.Errorf("longitudes don't match, expected %f got %f", testLong, location.GetLongitude())
	}
}
