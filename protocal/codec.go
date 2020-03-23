package protocal

import (
	"fmt"
	"strings"
)

type codec interface {
	Decode(in []byte) (Weight, error)
	GetDelimit() byte
}

type Weight struct {
	Value  float64
	Sign   byte
	Digits int
}

func (w Weight) String() string {
	var result []byte
	if MINUS == w.Sign {
		result = append(result, MINUS)
	}
	core := fmt.Sprintf("%v", w.Value)
	result = append(result, core...)
	var digitValue int
	if dotPos := strings.Index(core, "."); dotPos >= 0 {
		digitValue = len(core) - dotPos - 1
	}
	if digitDiff := w.Digits - digitValue; digitDiff > 0 {
		if digitValue == 0 {
			result = append(result, '.')
		}
		result = append(result, strings.Repeat("0", digitDiff)...)
	}
	return string(result)
}

func NewCodec(tf int) codec {
	switch tf {
	case 0:
		return tf0{}
	case 2:
		return tf23{
			frameLen: FRAME_LEN_TF2,
		}
	case 3:
		return tf23{
			frameLen: FRAME_LEN_TF3,
		}
	default:
		return tf0{}
	}
}
