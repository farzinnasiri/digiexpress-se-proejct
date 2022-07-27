package service

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/digiexpress/dlocator/internal/pkg/config"
	"github.com/digiexpress/dlocator/internal/pkg/log"
	"github.com/go-redis/redis/v9"
)

type CacheServiceImpl struct {
	redisClient *redis.Client
	ctx         *context.Context
	TTL         time.Duration
}

const ExpirationSet = "expiration"

func NewCacheService(config *cfg.AppConfig) (CacheServiceImpl, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(config.Redis.Host, ":", config.Redis.Port),
		Password: config.Redis.Password,
	})
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return CacheServiceImpl{}, err
	}

	TTL := time.Second * time.Duration(config.Redis.TTL)
	if TTL == 0 {
		log.Warning("TTL is set to 0")
	}

	return CacheServiceImpl{redisClient, &ctx, TTL}, nil
}

// AddGeoLocation O(log(1)) + O(log(N))
// where in worst case N := maximum number of active couriers
// more info: https://redis.io/commands/geoadd/
func (c CacheServiceImpl) AddGeoLocation(key string, courierDTO CourierDTO) error {
	if err := c.redisClient.GeoAdd(*c.ctx, key, &redis.GeoLocation{
		Name:      courierDTO.DriverId,
		Longitude: courierDTO.Longitude,
		Latitude:  courierDTO.Latitude,
	}).Err(); err != nil {
		return err
	}

	if err := c.updateExpirationSet(courierDTO.DriverId); err != nil {
		return err
	}

	return nil
}

// GeoSearchByRadius O(N+log(M)) + Housekeeping where N is the number of elements inside the bounding box of the circular area
//	delimited by center and radius and M is the number of items inside the index.
// more info:  https://redis.io/commands/georadius/
func (c CacheServiceImpl) GeoSearchByRadius(key string,
	queryByRadius SearchByRadiusDTO,
) ([]CourierDTO, error) {
	if err := c.removeDeadEntries(key); err != nil {
		return nil, err
	}

	locations, err := c.redisClient.GeoRadius(*c.ctx, key,
		queryByRadius.CenterLongitude, queryByRadius.CenterLatitude,
		&redis.GeoRadiusQuery{
			Radius:    queryByRadius.Radius,
			Unit:      "m",
			WithCoord: true,
			WithDist:  true,
			Sort:      "ASC",
		}).Result()
	if err != nil {
		return nil, err
	}

	couriers := make([]CourierDTO, len(locations))
	for i, location := range locations {
		couriers[i] = CourierDTO{
			DriverId:           location.Name,
			Latitude:           location.Latitude,
			Longitude:          location.Longitude,
			DistanceFromCenter: location.Dist,
		}
	}

	return couriers, nil
}

func (c CacheServiceImpl) updateExpirationSet(id string) error {
	now := time.Now()
	ttl := now.Add(c.TTL).Unix()
	return c.redisClient.ZAdd(*c.ctx, ExpirationSet, redis.Z{
		Score:  float64(ttl),
		Member: id,
	}).Err()
}

// removeDeadEntries  O(log(N)+M) + O(M*log(N))
func (c CacheServiceImpl) removeDeadEntries(key string) error {
	ids, err := c.redisClient.ZRangeByScore(*c.ctx, ExpirationSet, &redis.ZRangeBy{
		Min: "",
		Max: fmt.Sprint(time.Now().Unix()),
	}).Result()
	if err != nil {
		return err
	}
	if len(ids) != 0 {
		return c.redisClient.ZRem(*c.ctx, key, ids).Err()
	}

	return nil
}
