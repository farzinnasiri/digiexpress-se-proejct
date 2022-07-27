package config

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	TTL      int // in seconds
}
