package iblt

import (
	"fmt"
	"github.com/dchest/siphash"
	"testing"
)

func TestAbc(t *testing.T) {
	a := siphash.Hash(629, 465, []byte("abc"))
	fmt.Println(a)
	fmt.Println(a & 255)
	fmt.Println(byte(a))
}

func TestDataXor(t *testing.T) {
	d := make(data, 2)
	b := make(data, 2)
	d = []byte{2, 1}
	b = []byte{1, 2}
	d.xor(b)
	fmt.Println(d)
	fmt.Println(b)
}

func TestHashXor(t *testing.T) {
	a := hash(1)
	b := hash(2)
	a.xor(b)
	fmt.Println(a)
	fmt.Println(b)
}

func TestBucket_String(t *testing.T) {
	b := NewBucket(2)
	fmt.Println(b.empty())
	b.put([]byte{1, 3})
	fmt.Println(b)
	fmt.Println(b.pure())
	fmt.Println(b.empty())
	b.put([]byte{1, 4})
	fmt.Println(b)
	fmt.Println(b.pure())

	b1 := NewBucket(2)
	b1.put([]byte{1, 3})
	b.subtract(b1)
	fmt.Println(b)
	fmt.Println(b.pure())
}
