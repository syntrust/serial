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

func (w *ScaleReader) Listen(wt chan string, stop chan struct{}) error {
	c := &serial.Config{Name: w.PortName, Baud: w.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return fmt.Errorf("OpenPort failed: %v", err)
	}
	defer func() {
		if s != nil {
			if err = s.Close(); err != nil {
				fmt.Println("serial port close failed: ", err)
			}
			fmt.Println("serial port closed")
		}
		fmt.Println("defer", "s == nil")
	}()
	log.Println("connected to:", c.Name, "TF=", w.TF)

	reader := bufio.NewReader(s)
	var max, min, final float64
	timer := time.NewTimer(time.Second * time.Duration(w.Duration))

	readWeight := func() protocal.Weight {
		cdc := protocal.NewCodec(w.TF)
		raw, err := reader.ReadBytes(cdc.GetDelimit())
		if err != nil {
			log.Println("ReadBytes error", err)
		}
		w, err := cdc.Decode(raw)
		if err != nil {
			log.Println("Decode error", err)
		}
		return w
	}
	weight := readWeight()
	max = weight.Value
	min = weight.Value
	final = 0

	for {

		//log.Printf("read: %x=>%s", raw, weight.String())
		select {
		case <-timer.C:
			//after a few seconds of stable time, remember the maximum value ever
			final = max
			fmt.Println("set final", final)
		case <-stop:
			log.Println("Listen stopped")
			return nil
		default:
			weight := readWeight()
			if weight.Value > max {
				max = weight.Value
			} else if weight.Value < min {
				min = weight.Value
			}
			//it seems the truck is leaving when weight drops to 1/3 of the max
			if final > 0 && weight.Value < final/3 {
				theWeight := protocal.Weight{
					Value:  final,
					Sign:   weight.Sign,
					Digits: weight.Digits,
				}
				wt <- theWeight.String()
				fmt.Println("weigh success:", theWeight.String())
				return nil
			}
			//fmt.Println("max", max, "min", min)
			if final < 0 && max-min > float64(w.Deviation)/1000 {
				d := time.Second * time.Duration(w.Duration)
				timer.Reset(d)
				max = weight.Value
				min = weight.Value
				fmt.Println("Reset: max", max)
			}
		}
	}
}
