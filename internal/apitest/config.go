package apitest

import "github.com/spf13/viper"

func loadConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	return viper.ReadInConfig()
}
