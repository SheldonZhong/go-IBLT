package iblt

import (
	"fmt"
	"github.com/willf/bitset"
	"testing"
)

func TestBitSet(t *testing.T) {
	bts := bitset.New(20)
	bts.Set(0)
	fmt.Println(bts.Len())
	bts.Set(21)
	fmt.Println(bts.Len())
	fmt.Println(bts.Any())
	bts.ClearAll()
	fmt.Println(bts.Any())
}
