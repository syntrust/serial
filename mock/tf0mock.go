package mock

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type tf0 struct {
}

func (t tf0) send() {
	in := make(chan []byte)
	go serialOut(in, "/dev/ttyS1")
	//max width: 7; max precision: 4
	data := []float64{9.1, +20.00, 4252.97, 124.8209, 99984, 89.01, -1.8, 0.0, .622}
	fmts := []string{"%.1f", "%+.1f", "%7.2f", "%6.f", "%7.4f"}
	for _, f := range fmts {
		fmt.Println("=======" + f + "===========")
		for _, d := range data {
			formatted := fmt.Sprintf(f, d)
			item, err := t.encode(formatted)
			if err != nil {
				fmt.Println(err)
				continue
			}
			in <- item
			log.Printf("%s\t -> %x\n", formatted, item)
		}
	}
}

func (t tf0) encode(input string) ([]byte, error) {
	size, withDot := len(input), strings.Contains(input, ".")
	withSign := strings.Contains(input, string([]byte{PLUS})) || strings.Contains(input, string([]byte{MINUS}))
	if withSign {
		size -= 1
	}
	if withDot && size > 7 || !withDot && size > 6 {
		return nil, fmt.Errorf("input overflow: %s", input)
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return nil, err
	}
	encoded := make([]byte, FRAME_LEN)
	encoded[0] = STX
	encoded[11] = ETX
	if !math.Signbit(value) {
		encoded[1] = PLUS
	} else {
		encoded[1] = MINUS
		value = math.Abs(value)
	}
	var fracDig int
	if dotPos := strings.Index(input, "."); dotPos >= 0 {
		fracDig = len(input) - dotPos - 1
	}
	encoded[8] = byte(fracDig + OFFSET)
	numberIntTemp, _ := math.Modf(value * math.Pow10(fracDig))
	numberInt := int(numberIntTemp)
	for i := 7; i > 1; i-- {
		encoded[i] = byte(numberInt%10 + OFFSET)
		numberInt /= 10
	}
	encoded[9], encoded[10] = getXOR(encoded)
	return encoded, nil
}

func getXOR(encoded []byte) (h, l byte) {
	xor := 0
	for _, e := range encoded[1:9] {
		xor ^= int(e)
	}

	xorh := xor >> 4
	if xorh <= 9 {
		h = byte(xorh + OFFSET)
	} else {
		h = byte(xorh + X_OFFSET)
	}
	xorl := xor & 0xf
	if xorl <= 9 {
		l = byte(xorl + OFFSET)
	} else {
		l = byte(xorl + X_OFFSET)
	}
	return
}
