package mock

import (
	"github.com/tarm/serial"
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

func main() {
	tf23mock{FRAME_LEN_TF3}.send()
}

func serialOut(inputs chan []byte, portName string) {
	c := &serial.Config{Name: portName, Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	for item := range inputs {
		if _, err = s.Write(item); err != nil {
			panic(err)
		}
		//log.Printf("sent  %x\n", item)
	}
}
