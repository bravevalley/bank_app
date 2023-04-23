package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config is the struct containing the env to be exported
type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	Address       string        `mapstructure:"ADDRESS"`
	SymmetricKey  string        `mapstructure:"SYMMETRICKEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

// LoadConfig loads the config values to the Config struct and returns it
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// 	log.Fatal("Configuration file not found")
		// }
		// log.Fatal("Could not read config fge")
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		// if _, ok := err.(viper.ConfigMarshalError); ok {
		// 	log.Fatal("Unmarshall failed")
		// }
		// log.Fatal("Could not marshall the env")
		return
	}

	return
}
