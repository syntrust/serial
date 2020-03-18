package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"serialdemo/service"
)

var publicKey, _ = service.LoadPublicKey()

func scale(w http.ResponseWriter, r *http.Request) {
	var info service.WeightInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !verifySignature(info) {
		log.Println("invalid signature!")
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	// Do something with the WeightInfo struct...
	fmt.Fprintf(w, "WeightInfo: %+v", info)
	log.Printf("WeightInfo: %+v", info)
}

func verifySignature(info service.WeightInfo) bool {
	signInfo := &service.WeightInfo{
		Weight:    info.Weight,
		Vehicle:   info.Vehicle,
		ScaleSN:   info.ScaleSN,
		Location:  info.Location,
		TimeStamp: info.TimeStamp,
	}
	jsonBytes, _ := json.Marshal(signInfo)
	r, s := new(big.Int).SetBytes(info.R), new(big.Int).SetBytes(info.S)
	return ecdsa.Verify(publicKey, service.Hash(jsonBytes), r, s)
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
