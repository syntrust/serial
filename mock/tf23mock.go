package mock

import (
	"fmt"
	"log"
	"strings"
)

type tf23mock struct {
	frameLen int
}

func (t tf23mock) send() {
	in := make(chan []byte)
	go serialOut(in, "/dev/ttyS1")
	//max width: 7
	data := []float64{9.1, +20.00, 4252.97, 124.8209, 99984, 89.01, -1.8, 0.0, .622}
	for _, d := range data {
		formatted := fmt.Sprintf("%v", d)
		item, err := t.encode(formatted)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var stream []byte
		for i := 0; i < 5; i++ {
			stream = append(stream, item...)
			stream = append(stream, '=')
		}
		in <- stream
		log.Printf("%s\t -> %x\n", formatted, stream)
	}
}

func (t tf23mock) encode(input string) ([]byte, error) {
	if (strings.Index(input, ".") < 0 || strings.Index(input, ".") > t.frameLen-2) &&
		len(input) > t.frameLen-1 {
		return nil, fmt.Errorf("input overflow: %s", input)
	}
	if len(input) > 7 {
		input = input[:t.frameLen-1]
	}
	var result []byte
	for i := len(input) - 1; i >= 0; i-- {
		result = append(result, input[i])
	}
	for i := len(result); i < t.frameLen-1; i++ {
		result = append(result, '0')
	}
	return result, nil
}
