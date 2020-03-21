package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Post(infoToSign WeightInfoToSign, backendURL string) error {

	weightInfo := &WeightInfo{
		WeightInfoToSign: infoToSign,
	}
	var err error
	weightInfo.R, weightInfo.S, err = sign(infoToSign)
	if err != nil {
		return err
	}
	jsonValue, err := json.Marshal(weightInfo)
	if err != nil {
		return err
	}
	log.Println("posted", string(jsonValue))
	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	return nil
}
