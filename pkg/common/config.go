package common

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Configuration holds the data required by the application to run.
type Configuration struct {
	Version    string `mapstructure:"VERSION"`
	Port       int    `mapstructure:"PORT"`
	AppName    string `mapstructure:"APP_NAME"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	SSLMode    string `mapstructure:"DB_SSLMODE"`
}

// GetServerPort ...
func (c Configuration) GetServerPort() string {
	return fmt.Sprintf(":%d", c.Port)
}

// ReadConfig function read application configuration from (default) app.env file.
// This function does not return error, will throw Fatal instead.
func ReadConfig() Configuration {
	var config Configuration

	log.Info("reading config from app.env")
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed viper.ReadInConfig: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed unmarshal config: %v", err)
	}

	log.Info("success reading config from app.env")
	return config
}
