package service

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
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
	url           = "http://localhost:8080/scale"
)

type WeightReader struct {
	portName string
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
		ScaleSN:   "scale1",
		Location:  "location1",
		TimeStamp: time.Now().Unix(),
	}
	weightInfo, err := w.sign(infoToSign)
	if err != nil {
		return err
	}
	jsonValue, err := json.Marshal(weightInfo)
	if err != nil {
		return err
	}
	fmt.Println("post ", string(jsonValue))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	fmt.Println("response Status:", resp.Status)
	return nil
}

func (w *WeightReader) sign(infoToSign *WeightInfoToSign) (*WeightInfo, error) {
	jsonValue, _ := json.Marshal(infoToSign)
	fmt.Println("sign ", string(jsonValue))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, Hash(jsonValue))
	if err != nil {
		return nil, err
	}
	info := &WeightInfo{
		WeightInfoToSign: *infoToSign,
	}
	info.R, info.S = r.Bytes(), s.Bytes()
	return info, nil
}

func Hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}
