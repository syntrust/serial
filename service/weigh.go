package service

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"log"
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
	Timeout    int
	Deviation  int
	BackendURL string
}

type ScaleReader struct {
	*Config
}

type WeightInfo struct {
	Error string
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

func (w *ScaleReader) Listen(result chan string, stop chan struct{}) error {
	c := &serial.Config{Name: w.PortName, Baud: w.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("OpenPort failed: %v", err)
	}
	reader := bufio.NewReader(s)
	defer func() {
		if s != nil {
			_ = s.Close()
		}
	}()
	log.Println("connected to:", c.Name, "TF=", w.TF)
	wCh := make(chan protocal.Weight)

	go func() {
		cdc := protocal.NewCodec(w.TF)
		for {
			raw, err := reader.ReadBytes(cdc.GetDelimit())
			if err != nil {
				log.Println("ReadBytes error", err)
				//continue
			}
			w, err := cdc.Decode(raw)
			if err != nil {
				log.Println("Decode error", err)
				//continue
			}
			wCh <- w
			//fmt.Println("for weight", w)
		}
	}()
	weight := <-wCh
	fmt.Println("first weight", weight)
	maxEver, max, min := weight.Value, weight.Value, weight.Value
	timer := time.NewTimer(time.Second * time.Duration(w.Duration))
	var final float64 = 0
	for {
		select {
		case <-timer.C:
			//remember the maximum value during last stable time window
			final = max
			fmt.Println("set final", final)
		case <-stop:
			log.Println("Listen timeout, will stop. max value ever: ", maxEver)
			return nil
		case weight = <-wCh:
			if weight.Value > max {
				max = weight.Value
				//fmt.Println("set max=", max)
			} else if weight.Value < min {
				min = weight.Value
				//fmt.Println("set min=", min)
			}
			//it seems the truck is leaving when weight drops to 1/3 of the max
			if final > 0 && weight.Value < final/3 {
				theWeight := protocal.Weight{
					Value:  final,
					Sign:   weight.Sign,
					Digits: weight.Digits,
				}
				result <- theWeight.String()
				fmt.Println("weigh success:", theWeight.String())
				return nil
			}
			//fmt.Println("max", max, "min", min)
			if final == 0 && max-min > float64(w.Deviation)/1000 {
				d := time.Second * time.Duration(w.Duration)
				timer.Reset(d)
				if maxEver < max {
					maxEver = max
				}
				max = weight.Value
				min = weight.Value
				log.Println("Reset: max", max)
			}
		default:

		}
	}
}
