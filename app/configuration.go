package app

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/danielkraic/idmapper/app/idmappers"
)

// Configuration application configuration
type Configuration struct {
	Addr       string           `mapstructure:"addr"`
	APIPrefix  string           `mapstructure:"api_prefix"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Redis      RedisConfig      `mapstructure:"redis"`
	PostgreSQL PostgreSQLConfig `mapstructure:"pgsql"`
	IDMappers  idmappers.Config `mapstructure:"idmappers"`
}

// LoggerConfig application configuration for Logger
type LoggerConfig struct {
	JSON bool `mapstructure:"json"`
}

// RedisConfig application configuration for Redis
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"pasword"`
}

// PostgreSQLConfig application configuration for PostgreSQL
type PostgreSQLConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

func readConfiguration(configFile string) (*Configuration, error) {
	viper.SetConfigFile(configFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("IDMAPPER")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, fmt.Errorf("unable to BindPFlags, %s", err)
	}

	viper.SetDefault("addr", "0.0.0.0:80")
	viper.SetDefault("api_prefix", "/v1")
	viper.SetDefault("logger.json", false)
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("idmappers.reloader.intervals.currency", "24h")
	viper.SetDefault("idmappers.reloader.intervals.country", "24h")
	viper.SetDefault("idmappers.reloader.intervals.language", "24h")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			return nil, fmt.Errorf("failed to read config file: %s", err)
		}
	}

	var configuration Configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode configration to struct, %s", err)
	}

	return &configuration, nil
}
