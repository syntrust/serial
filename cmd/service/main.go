package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"serialdemo/service"
	"time"
)

var (
	timeout = time.Second * 10
	conf    service.Config
)

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

func barOpen(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	v := vars.Get("vehicle")
	checkout := vars.Get("checkout") == "true"
	log.Println("drive in:", v, "checkout", checkout)
	weightChan := make(chan string)
	go func() {
		if err := service.NewWeightReader(&conf).Listen(weightChan); err != nil {
			log.Fatal(err)
		}
	}()
	weit := <-weightChan
	infoToSign := service.WeightInfoToSign{
		Weight:    weit,
		Vehicle:   v,
		ScaleSN:   conf.ScaleSN,
		Location:  conf.Location,
		Checkout:  checkout,
		TimeStamp: time.Now().Unix(),
	}
	if err := service.Post(infoToSign, conf.BackendURL); err != nil {
		log.Println("ERROR", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bar", barOpen)
	url := "0.0.0.0:9090"
	log.Println("checkpoint listening:", url)
	if err := http.ListenAndServe(url, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
