package protocal

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strings"
)

type codec interface {
	encode(in string) ([]byte, error)
	decode(in []byte) (weight, error)
	getDelimit() byte
}

type weight struct {
	value  float64
	sign   byte
	digits int
}

func (w weight) String() string {
	var result []byte
	if MINUS == w.sign {
		result = append(result, MINUS)
	}
	core := fmt.Sprintf("%v", w.value)
	result = append(result, core...)
	var digitValue int
	if dotPos := strings.Index(core, "."); dotPos >= 0 {
		digitValue = len(core) - dotPos - 1
	}
	if digitDiff := w.digits - digitValue; digitDiff > 0 {
		if digitValue == 0 {
			result = append(result, '.')
		}
		result = append(result, strings.Repeat("0", digitDiff)...)
	}
	return string(result)
}

type WeightReader struct {
	portName string
}

func NewWeightReader(portName string) WeightReader {
	return WeightReader{portName: portName}
}

func (w *WeightReader) Listen(tf int) {
	c := &serial.Config{Name: w.portName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	log.Println("connected:", c.Name)
	reader := bufio.NewReader(s)
	cdc := NewCodec(tf)
	log.Println("listening to: TF=", tf)
	for {
		source, err := reader.ReadBytes(cdc.getDelimit())
		if err != nil {
			panic(err)
		}
		weight, err := cdc.decode(source)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%x=>%s", source, weight.String())
	}
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
