package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Post(msg interface{}, backendURL string) error {
	var jsonValue []byte
	var err error
	switch value := msg.(type) {
	case string:
		weightInfo := &WeightInfo{
			Error: value,
		}
		jsonValue, err = json.Marshal(weightInfo)
		if err != nil {
			return err
		}
		log.Println("case string", string(jsonValue))
	case WeightInfoToSign:
		weightInfo := &WeightInfo{
			WeightInfoToSign: value,
		}
		weightInfo.R, weightInfo.S, err = sign(value)
		if err != nil {
			return err
		}
		jsonValue, err = json.Marshal(weightInfo)
		if err != nil {
			return err
		}
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
