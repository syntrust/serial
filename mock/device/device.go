package device

import (
	"github.com/tarm/serial"
	"time"
)

const (
	FRAME_LEN     = 12
	FRAME_LEN_TF2 = 8
	FRAME_LEN_TF3 = 9
	STX           = '\x02'
	ETX           = '\x03'
	PLUS          = '\x2B'
	MINUS         = '\x2D'
	OFFSET        = 0x30
	X_OFFSET      = 0x37
)

type Weigh interface {
	encode(input string) ([]byte, error)
	Send(data float64, in chan []byte)
}

func NewMock(tf int) Weigh {
	switch tf {
	case 0:
		return Tf0Mock{}
	case 2:
		return Tf23Mock{
			frameLen: FRAME_LEN_TF2,
		}
	case 3:
		return Tf23Mock{
			frameLen: FRAME_LEN_TF3,
		}
	default:
		return Tf0Mock{}
	}
}
func SerialOut(inputs chan []byte, portName string) {
	c := &serial.Config{Name: portName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	for item := range inputs {
		if _, err = s.Write(item); err != nil {
			panic(err)
		}
		//paud=9600
		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}
