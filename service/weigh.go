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
	SiteSN     string
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
	SiteSN    string
	Checkout  bool
	TimeStamp int64
}

func (w *ScaleReader) Listen(result chan string, stop chan struct{}) error {

	wCh := make(chan protocal.Weight)
	errCh := make(chan error)
	quitCh := make(chan struct{})
	go func() {
		c := &serial.Config{Name: w.PortName, Baud: w.Baud}
		s, err := serial.OpenPort(c)
		if err != nil {
			errCh <- fmt.Errorf("OpenPort failed: %v", err)
			return
		}
		defer func() {
			if s != nil {
				err = s.Close()
				log.Println("s.Close()),", err)
			} else {

				log.Println("s  == nil ),", err)
			}
		}()
		log.Println("connected to:", c.Name, "TF=", w.TF)
		cdc := protocal.NewCodec(w.TF)
		reader := bufio.NewReader(s)
		go func() {
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
				//fmt.Println("for weight", w)
				wCh <- w
			}
		}()

		<-quitCh
		log.Println("loop quit")
		return
	}()
	var weight protocal.Weight
	select {
	case weight = <-wCh:
		fmt.Println("first weight", weight)
	case err := <-errCh:
		return err
	}
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
			quitCh <- struct{}{}
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
				log.Println("truck is leaving...")
				go func() {
					quitCh <- struct{}{}
				}()
				result <- theWeight.String()
				log.Println("weigh success:", theWeight.String())
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
