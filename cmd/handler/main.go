package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serialdemo/service"
)

func scale(w http.ResponseWriter, r *http.Request) {
	var info service.WeightInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Do something with the WeightInfo struct...
	fmt.Fprintf(w, "WeightInfo: %+v", info)
	log.Printf("WeightInfo: %+v", info)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/scale", scale)
	url := "0.0.0.0:8080"
	if err := http.ListenAndServe(url, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("ListenAndServe", url)
}
