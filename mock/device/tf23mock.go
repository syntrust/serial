package device

import (
	"fmt"
	"log"
	"strings"
)

type Tf23Mock struct {
	frameLen int
}

func (t Tf23Mock) Send(data float64, in chan []byte) {
	//max width: 7
	formatted := fmt.Sprintf("%v", data)
	item, err := t.encode(formatted)
	if err != nil {
		fmt.Println(err)
		return
	}
	item = append(item, '=')
	in <- item
	log.Printf("%s\t -> %x\n", formatted, item)
}

func (t Tf23Mock) encode(input string) ([]byte, error) {
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
