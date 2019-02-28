package iblt

import (
	"fmt"
	"github.com/dchest/siphash"
)

const (
	key0 = 465
	key1 = 629
)

func sipHash(b []byte) uint64 {
	// TODO: key constants
	return siphash.Hash(key0, key1, b)
}

type data []byte
type hash byte

func xor(dst []byte, src []byte) {
	for i, v := range dst {
		dst[i] = v ^ src[i]
	}
}

func (d data) xor(a data) {
	xor(d, a)
}

func (d data) empty() bool {
	for _, v := range d {
		if v != byte(0) {
			return false
		}
	}

	return true
}

func (d *hash) xor(a hash) {
	*d = *d ^ a
}

func (d hash) empty() bool {
	return d == hash(0)
}

type Bucket struct {
	dataSum data
	hashSum hash
	count   int
}

func NewBucket(len int) *Bucket {
	return &Bucket{
		dataSum: make(data, len),
		hashSum: hash(0),
		count:   0,
	}
}

func (b *Bucket) xor(a *Bucket) {
	b.dataSum.xor(a.dataSum)
	b.hashSum.xor(a.hashSum)
}

func (b *Bucket) add(a *Bucket) {
	b.xor(a)
	b.count = b.count + a.count
}

func (b *Bucket) subtract(a *Bucket) {
	b.xor(a)
	b.count = b.count - a.count
}

func (b *Bucket) put(d data) {
	b.dataSum.xor(d)
	h := sipHash(d)
	b.hashSum.xor(hash(h))
	b.count++
}

func (b *Bucket) take(d data) {
	b.dataSum.xor(d)
	h := sipHash(d)
	b.hashSum.xor(hash(h))
	b.count--
}

func (b Bucket) copy() *Bucket {
	bkt := NewBucket(len(b.dataSum))
	copy(bkt.dataSum, b.dataSum)
	bkt.hashSum = b.hashSum
	bkt.count = b.count
	return bkt
}

func (b Bucket) pure() bool {
	if b.count == 1 || b.count == -1 {
		h := sipHash(b.dataSum)
		fmt.Printf("pure: hash: %v, data: %v\n", hash(h), b.dataSum)
		if b.hashSum == hash(h) {
			return true
		}
	}
	return false
}

func (b Bucket) empty() bool {
	return b.count == 0 &&
		b.hashSum.empty() &&
		b.dataSum.empty()
}

func (b Bucket) String() string {
	return fmt.Sprintf("Bucket: dataSum: %v, hashSum: %v, count: %d",
		b.dataSum, b.hashSum, b.count)
}

// each part of symmetric difference
type Diff struct {
	Alpha [][]byte
	Beta  [][]byte
}

func NewDiff() *Diff {
	return &Diff{}
}

// assume b is pure bucket
func (d *Diff) encode(b *Bucket) {
	cpy := make([]byte, len(b.dataSum))
	copy(cpy, b.dataSum)
	if b.count == 1 {
		d.Alpha = append(d.Alpha, cpy)
	}
	if b.count == -1 {
		d.Beta = append(d.Beta, cpy)
	}
}

func (d Diff) String() string {
	s := "Symmetric Difference:\n"
	for i, b := range d.Alpha {
		s += fmt.Sprintf("Alpha:%d:%v\n", i, b)
	}

	for i, b := range d.Beta {
		s += fmt.Sprintf("Beta:%d:%v\n", i, b)
	}

	return s
}
