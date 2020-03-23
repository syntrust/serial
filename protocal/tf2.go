package protocal

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	EQUAL         = '\x3D'
	FRAME_LEN_TF2 = 8
	FRAME_LEN_TF3 = 9
)

type tf23 struct {
	frameLen int
}

func (t tf23) GetDelimit() byte {
	return EQUAL
}
func (t tf23) Decode(source []byte) (Weight, error) {
	if len(source) != t.frameLen || source[t.frameLen-1] != EQUAL {
		return Weight{}, fmt.Errorf("invalid data: %x", source)
	}
	var vb []byte
	sign := PLUS
	for i := t.frameLen - 2; i >= 0; i-- {
		if i == t.frameLen-2 && source[i] == MINUS {
			sign = MINUS
			for ; source[i] == '0'; i-- {
				continue
			}
		}
		vb = append(vb, byte(source[i]))
	}
	vb = bytes.TrimRight(vb, ".")
	vb = bytes.TrimLeft(vb, "0")
	if len(vb) == 0 {
		vb = []byte{'0'}
	}
	v, err := strconv.ParseFloat(string(vb), 64)
	if err != nil {
		return Weight{}, err
	}
	return Weight{
		Value:  v,
		Sign:   byte(sign),
		Digits: bytes.Index(source, []byte(".")),
	}, nil
}
