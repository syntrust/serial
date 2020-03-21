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
	cameraURL = "http://localhost:9090/bar"
	duration  = 5
)

var data = []float64{9.1, +20.00, 4252.97, 124.8209, 99984, 89.01, -1.8, 0.0, .622}

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
		r := rand.Intn(len(data))
		timer := time.NewTimer(time.Second * duration)
	out:
		for {
			mck.Send(data[r], in)
			select {
			case <-timer.C:
				break out
			default:
				continue
			}
		}
	}()
	r := <-res
	fmt.Println("response Status:", r)
}
