package config

import "github.com/spf13/viper"

func Init() error {
	viper.AddConfigPath("./config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
