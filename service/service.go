package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"net/http"
	"serialdemo/protocal"
	"time"
)

var (
	privateKey, _ = loadPrivateKey("private.key")
)

type Config struct {
	ScaleSN    string
	Location   string
	PortName   string
	Baud       int
	TF         int
	Duration   int
	Interval   int
	Deviation  int
	BackendURL string
}

type WeightReader struct {
	*Config
}

type WeightInfo struct {
	WeightInfoToSign
	R []byte
	S []byte
}
type WeightInfoToSign struct {
	Weight    string
	Vehicle   string
	ScaleSN   string
	Location  string
	TimeStamp int64
}

func NewWeightReader(conf *Config) *WeightReader {
	return &WeightReader{
		conf,
	}
}

func (w *WeightReader) Listen() error {
	c := &serial.Config{Name: w.PortName, Baud: w.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	log.Println("connected to :", c.Name)
	reader := bufio.NewReader(s)
	cdc := protocal.NewCodec(w.TF)
	log.Println("listening to: TF=", w.TF)
	for {
		source, err := reader.ReadBytes(cdc.GetDelimit())
		if err != nil {
			log.Println("ERROR", err)
		}
		weight, err := cdc.Decode(source)
		if err != nil {
			log.Println("ERROR", err)
		}
		log.Printf("%x=>%s", source, weight.String())
		err = w.post(weight.String())
		if err != nil {
			log.Println("ERROR", err)
		}
	}
}

func (w *WeightReader) post(weight string) error {
	infoToSign := &WeightInfoToSign{
		Weight:    weight,
		Vehicle:   "vehicle1",
		ScaleSN:   w.ScaleSN,
		Location:  w.Location,
		TimeStamp: time.Now().Unix(),
	}
	weightInfo := &WeightInfo{
		WeightInfoToSign: *infoToSign,
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
	fmt.Println("post ", string(jsonValue))
	resp, err := http.Post(w.BackendURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	fmt.Println("response Status:", resp.Status)
	return nil
}
