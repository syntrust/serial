package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"serialdemo/protocal"
	"time"
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

func Weigh(ctx context.Context, conf Config, wCh <-chan protocal.Weight, result chan<- string) error {

	defer func() {
		log.Println("s.Close()")
	}()

	var weight protocal.Weight
	weight = <-wCh
	fmt.Println("first weight", weight)

	maxEver, max, min := weight.Value, weight.Value, weight.Value
	timer := time.NewTimer(time.Second * time.Duration(conf.Duration))
	var final float64 = 0
	for {
		select {
		case <-timer.C:
			//remember the maximum value during last stable time window
			final = max
			fmt.Println("set final", final)
		case <-ctx.Done():
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
				result <- theWeight.String()
				log.Println("weigh success:", theWeight.String())
				return nil
			}
			//fmt.Println("max", max, "min", min)
			if final == 0 && max-min > float64(conf.Deviation)/1000 {
				d := time.Second * time.Duration(conf.Duration)
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

func ListenWeight(conf Config, wCh chan protocal.Weight, quitCh chan struct{}) error {
	c := &serial.Config{Name: conf.PortName, Baud: conf.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("OpenPort failed: %v", err)
	}
	defer func() {
		_ = s.Close()
		log.Println("s.Close()")
	}()
	log.Println("connected to:", c.Name, "TF=", conf.TF)

	reader := bufio.NewReader(s)
	cdc := protocal.NewCodec(conf.TF)
	for {
		select {
		case <-quitCh:
			//reader.Reset(s)
			//log.Println("reader.Reset()")
			return nil
		default:
			raw, err := reader.ReadBytes(cdc.GetDelimit())
			if err != nil {
				log.Println("ReadBytes error", err)
				//continue
			}
			decoded, err := cdc.Decode(raw)
			if err != nil {
				log.Println("Decode error", err)
				//continue
			}

			wCh <- decoded
		}
	}
}
