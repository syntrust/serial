package main

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"serialdemo/service"
	"time"
)

var (
	conf service.Config
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
	log.Println("truck", v, "checkout", checkout)
	weightChan := make(chan string)
	errChan := make(chan error)
	stopChan := make(chan struct{})
	var msg interface{}
	go func() {
		read := service.ScaleReader{
			Config: &conf,
		}
		if err := read.Listen(weightChan, stopChan); err != nil {
			log.Println("ERROR", err)
			errChan <- err
		}
	}()
	select {
	case weit := <-weightChan:
		msg = service.WeightInfoToSign{
			Weight:    weit,
			Vehicle:   v,
			ScaleSN:   conf.ScaleSN,
			Location:  conf.Location,
			Checkout:  checkout,
			TimeStamp: time.Now().Unix(),
		}
	case err := <-errChan:
		msg = err.Error()
		log.Println("ERROR", err)
	case <-time.After(time.Second * time.Duration(conf.Timeout)):
		stopChan <- struct{}{}
		msg = "ScaleReader timeout"
		log.Println("ERROR", msg)
	}
	if err := service.Post(msg, conf.BackendURL); err != nil {
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
