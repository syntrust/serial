package main

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	var msg interface{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(conf.Timeout))
	defer cancel()
	go func() {
		sr := service.ScaleReader{
			Config: &conf,
		}
		if err := sr.Listen(ctx, weightChan); err != nil {
			errChan <- err
		}
	}()
	select {
	case weit := <-weightChan:
		msg = service.WeightInfoToSign{
			Weight:    weit,
			Vehicle:   v,
			ScaleSN:   conf.ScaleSN,
			SiteSN:    conf.SiteSN,
			Checkout:  checkout,
			TimeStamp: time.Now().Unix(),
		}
	case err := <-errChan:
		msg = err.Error()
		log.Println("ERROR", err)
	case <-ctx.Done():
		msg = "ScaleReader timeout"
		log.Println("ERROR", msg)
	}
	if err := service.Post(msg, conf.BackendURL); err != nil {
		log.Println("ERROR", err)
	}
	//stopChan <- struct{}{}
	<-ctx.Done()
	msg = "process timeout"
	log.Println("ERROR", msg)
}

func main() {
	go func() {

		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/bar", barOpen)
	url := "0.0.0.0:9090"
	log.Println("checkpoint listening:", url)
	if err := http.ListenAndServe(url, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
