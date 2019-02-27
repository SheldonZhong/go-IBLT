package iblt

import (
	"fmt"
	"github.com/dchest/siphash"
	"testing"
)

var testKey = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

func TestIBLT(t *testing.T) {
	h := siphash.New(testKey)
	fmt.Println(h.Write([]byte{1, 2, 3, 4, 5}))
	fmt.Println(h.Sum(nil))

	table := NewTable(40, 8, 4)
	err := table.Insert(h.Sum(nil))
	if err != nil {
		fmt.Println("err", err)
	}
	d, err := table.Decode()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(d)
}
