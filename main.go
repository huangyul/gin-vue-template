package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	initViper()
	s := InitServer()

	s.Run(fmt.Sprintf("127.0.0.1:%d", viper.GetInt("server.port")))
}

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
