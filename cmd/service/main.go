package main

import (
	"flag"
	"serialdemo/service"
)

func main() {
	tf := flag.Int("tf", 0, "指定通讯方式（0~8）")
	flag.Parse()
	pr := service.NewWeightReader("COM1")
	pr.Listen(*tf)
}
