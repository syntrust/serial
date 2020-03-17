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

type WeightReader struct {
	portName string
}

func NewWeightReader(portName string) WeightReader {
	return WeightReader{portName: portName}
}

func (w *WeightReader) Listen(tf int) {
	c := &serial.Config{Name: w.portName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	log.Println("connected:", c.Name)
	reader := bufio.NewReader(s)
	cdc := protocal.NewCodec(tf)
	log.Println("listening to: TF=", tf)
	for {
		source, err := reader.ReadBytes(cdc.GetDelimit())
		if err != nil {
			panic(err)
		}
		weight, err := cdc.Decode(source)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%x=>%s", source, weight.String())
		post(weight.String())
	}
}

type WeightInfo struct {
	Weight    string
	Vehicle   string
	ScaleSN   string
	Location  string
	TimeStamp int64
}

func post(w string) {

	fmt.Println("posting weight", w)
	weightInfo := &WeightInfo{
		Weight:    w,
		Vehicle:   "vehicle1",
		ScaleSN:   "scale1",
		Location:  "location1",
		TimeStamp: time.Now().Unix(),
	}

	fmt.Println("posting weightInfo", weightInfo)
	jsonValue, _ := json.Marshal(weightInfo)

	fmt.Println("posting jsonValue", string(jsonValue))
	url := "http://localhost:8080/scale"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response Status:", resp.Status)
}
