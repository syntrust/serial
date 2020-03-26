package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"serialdemo/mock/device"
	"time"
)

const (
	cameraURL     = "http://localhost:9090/bar"
	duration      = 5
	desiredStdDev = 0.05
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
		fmt.Println("rand", r)
		timer := time.NewTimer(time.Second * duration)
		raw := data[r]

		for {
			sample := rand.NormFloat64()*desiredStdDev + raw
			mck.Send(sample, in)
			select {
			case <-timer.C:
				//mock truck leaving
				mck.Send(raw/2, in)
				time.Sleep(time.Second)
				mck.Send(raw/4, in)
				return
			default:
				continue
			}
		}
	}()
	r := <-res
	fmt.Println("response Status:", r)
}
