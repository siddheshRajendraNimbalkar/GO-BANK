package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver    string        `mapstructure:"DBDRIVE"`
	DBSource    string        `mapstructure:"DBSOURCE"`
	Addr        string        `mapstructure:"ADDR"`
	Secret      string        `mapstructure:"SECRET"`
	JwtDuration time.Duration `mapstructure:"JWTDURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
