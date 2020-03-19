package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"serialdemo/service"
)

func main() {
	var conf service.Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	if err := service.NewWeightReader(&conf).Listen(); err != nil {
		log.Fatal(err)
	}
}
