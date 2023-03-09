package config

import "github.com/spf13/viper"

type Config struct {
	DbDriver          string `mapstructure:"DB_DRIVER"`
	DbAddress         string `mapstructure:"DB_ADDRESS"`
	HttpServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func Parse(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
