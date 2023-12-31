package app

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"runtime"
)

type Env struct {
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	DBType            string `mapstructure:"DB_TYPE"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPass            string `mapstructure:"DB_PASS"`
	DBName            string `mapstructure:"DB_NAME"`
	OpenWeatherAPIKey string `mapstructure:"OW_API_KEY"`
}

func NewEnv(isTestEnv bool) *Env {
	env := Env{}
	if isTestEnv {
		_, filename, _, _ := runtime.Caller(1)
		dir := filepath.Dir(filename)
		viper.SetConfigFile(dir + "/../testing.env")
	} else {
		viper.SetConfigFile(".env")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
