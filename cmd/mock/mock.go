package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"serialdemo/mock/device"
	"time"
)

const (
	cameraURL     = "http://localhost:9090/bar"
	duration      = 5
	desiredStdDev = 0.001
)

var data = []float64{9.1, +20.00, 422.97, 124.8209, 984, 89.01, -1.8, .622}

func main() {
	//truck loaded and checkout
	url1, _ := url.Parse(cameraURL)
	params := url.Values{}
	params.Set("vehicle", "陕A888888")
	params.Set("checkout", "true")
	url1.RawQuery = params.Encode()
	res := make(chan string)
	go func() {
		resp, err := http.Get(url1.String())
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		res <- resp.Status
	}()
	tf := flag.Int("tf", 0, "指定通讯方式（0~8）")
	flag.Parse()

	go func() {
		mck := device.NewMock(*tf)
		in := make(chan []byte)
		go device.SerialOut(in, "COM2")
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(data))
		timer := time.NewTimer(time.Second * duration)
		timer1 := time.NewTimer(time.Millisecond * time.Duration(duration*1000-500))
		raw := data[r]
		total := 0
		for {
			select {
			case <-timer.C:
				mck.Send(data[r]/4, in)
				log.Println("total sent", total)
				return
			case <-timer1.C:
				//mock truck leaving
				raw = data[r] / 2
			default:
				mck.Send(rand.NormFloat64()*desiredStdDev+raw, in)
				total++
			}
		}
	}()
	r := <-res
	fmt.Println("response Status:", r)
}
