package iblt

import (
	"fmt"
	"github.com/willf/bitset"
	"testing"
)

func TestBitSet(t *testing.T) {
	bts := bitset.New(21)
	bts.Set(0)
	fmt.Println(bts.Len())
	bts.Set(15)
	fmt.Println(bts.Len())

	//for i := 0; bts.Any(); i ++ {
	//	if bts.Test(uint(i)) {
	//		fmt.Println(i)
	//		bts.Clear(uint(i))
	//	}
	//}
	i := uint(0)
	v := true
	for {
		i, v = bts.NextSet(i)
		fmt.Println(i, v)
		if v {
			bts.Clear(i)
		} else {
			break
		}
	}
}
