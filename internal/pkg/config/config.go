package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	AppName         = "dlocator"
	FilePath        = "./configs/config.yml"
	DefaultGrpcPort = 8080
)

type AppConfig struct {
	Logging LoggingConfig
	Grpc    GrpcConfig
}

func NewAppConfig() (*AppConfig, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configurations")
	}

	return config, nil
}

func setDefaultConfigs() {
	viper.SetDefault("log.level", "warn")
	viper.SetDefault("grpc.appName", AppName)
	viper.SetDefault("grpc.port", DefaultGrpcPort)
}

func readEnvConfigs() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(strings.ToUpper(AppName))
	viper.AutomaticEnv()
}

func readFileConfigs() error {
	viper.SetConfigFile(FilePath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fialed to read file config: %w", err)
	}

	return nil
}

func getAggregatedConfigs() (*AppConfig, error) {
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func loadConfig() (*AppConfig, error) {
	setDefaultConfigs()
	readEnvConfigs()

	if err := readFileConfigs(); err != nil {
		return nil, err
	}

	return getAggregatedConfigs()
}
