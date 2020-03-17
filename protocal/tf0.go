package protocal

import (
	"fmt"
	"strconv"
)

const (
	OFFSET   = 0x30
	X_OFFSET = 0x37
	STX      = '\x02'
	ETX      = '\x03'
	PLUS     = '\x2B'
	MINUS    = '\x2D'
)

type tf0 struct {
}

func (t tf0) GetDelimit() byte {
	return ETX
}

func (t tf0) Decode(source []byte) (weight, error) {
	if len(source) != 12 || source[0] != STX || source[11] != ETX {
		return weight{}, fmt.Errorf("invalid data: %x", source)
	}
	h, l := getXOR(source)
	if h != source[9] || l != source[10] {
		return weight{}, fmt.Errorf("xor validation failed: %x", source)
	}
	d := int(source[8] - OFFSET)
	var vb []byte
	vb = append(vb, source[2:8-d]...)
	vb = append(vb, '.')
	vb = append(vb, source[8-d:8]...)
	v, err := strconv.ParseFloat(string(vb), 64)
	if err != nil {
		return weight{}, err
	}
	return weight{
		value:  v,
		sign:   source[1],
		digits: d,
	}, nil
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
