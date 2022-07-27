package service

import (
	"context"
	"testing"

	"github.com/digiexpress/dlocator/internal/pkg/config"
)

func setupTest() CacheServiceImpl {
	cacheService, err := NewCacheService(&config.AppConfig{
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "1234",
			TTL:      300,
		},
	})
	if err != nil {
		panic(err)
	}
	cacheService.redisClient.FlushAll(context.Background())
	return cacheService
}

func tearDownTest(cacheService CacheServiceImpl) {
	cacheService.redisClient.FlushAll(context.Background())
}

func TestCacheServiceAddGeoLocation(t *testing.T) {
	cacheService := setupTest()

	locations := []CourierDTO{
		{
			DriverId:  "1",
			Latitude:  35,
			Longitude: 53,
		}, {
			DriverId:  "2",
			Latitude:  24,
			Longitude: 16,
		}, {
			DriverId:  "3",
			Latitude:  19,
			Longitude: 100,
		},
	}

	for _, location := range locations {
		if err := cacheService.AddGeoLocation("test", location); err != nil {
			t.Errorf(err.Error())
		}
	}

	values, err := cacheService.redisClient.ZRange(context.Background(), "test", 0, 180).Result()
	if err != nil {
		t.Errorf(err.Error())
	}

	for _, value := range values {
		if !(value == "1" || value == "2" || value == "3") {
			t.Errorf("value is not in the expected range")
		}
	}

	tearDownTest(cacheService)
}

func TestGeoSearchByRadius(t *testing.T) {
	cacheService := setupTest()
	locations := []CourierDTO{
		{
			DriverId:  "1",
			Latitude:  38.115556,
			Longitude: 13.361389,
		}, {
			DriverId:  "2",
			Latitude:  41.502669,
			Longitude: 15.087269,
		},
	}

	for _, location := range locations {
		if err := cacheService.AddGeoLocation("test", location); err != nil {
			t.Errorf(err.Error())
		}
	}

	locationDTOs, err := cacheService.GeoSearchByRadius("test", SearchByRadiusDTO{
		CenterLatitude:  37,
		CenterLongitude: 15,
		Radius:          200_000,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(locationDTOs) > 1 {
		t.Errorf("no more than a single result should be returend")
		t.FailNow()
	}

	if locationDTOs[0].DriverId != "1" {
		t.Errorf("incorrect id is returned")
	}

	tearDownTest(cacheService)
}
