package main

import (
	"flag"
	"serialdemo/mock/device"
)

func main() {
	tf := flag.Int("tf", 0, "指定通讯方式（0~8）")
	flag.Parse()
	mck := device.NewMock(*tf)
	in := make(chan []byte)
	go device.SerialOut(in, "/dev/ttyS1")
	data := []float64{9.1, +20.00, 4252.97, 124.8209, 99984, 89.01, -1.8, 0.0, .622}
	mck.Send(data, in)
}
