package service

import (
	"bufio"
	"github.com/tarm/serial"
	"log"
	"serialdemo/protocal"
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
	Checkout  bool
	TimeStamp int64
}

func NewWeightReader(conf *Config) *WeightReader {
	return &WeightReader{
		conf,
	}
}

func (w *WeightReader) Listen(wt chan string) error {
	c := &serial.Config{Name: w.PortName, Baud: w.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	log.Println("connected to :", c.Name, "TF=", w.TF)
	reader := bufio.NewReader(s)
	cdc := protocal.NewCodec(w.TF)
	for {
		raw, err := reader.ReadBytes(cdc.GetDelimit())
		if err != nil {
			log.Println("ERROR", err)
		}
		weight, err := cdc.Decode(raw)
		if err != nil {
			log.Println("ERROR", err)
		}
		wt <- weight.String()
		log.Printf("read: %x=>%s", raw, weight.String())
		return nil
	}
}
