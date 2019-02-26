package iblt

import (
	"fmt"
	"github.com/dchest/siphash"
)

func sipHash(b []byte) uint64 {
	// TODO: key constants
	return siphash.Hash(465, 629, b)
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

func NewBucket(len int) Bucket {
	return Bucket{
		dataSum: make(data, len),
		hashSum: hash(0),
		count:   0,
	}
}

func (b *Bucket) xor(a Bucket) {
	b.dataSum.xor(a.dataSum)
	b.hashSum.xor(a.hashSum)
}

func (b *Bucket) add(a Bucket) {
	b.xor(a)
	b.count = b.count + a.count
}

func (b *Bucket) subtract(a Bucket) {
	b.xor(a)
	b.count = b.count - a.count
}

func (b *Bucket) put(d data) {
	b.dataSum.xor(d)
	h := sipHash(d)
	b.hashSum.xor(hash(h))
	b.count++
}

func (b Bucket) pure() bool {
	if b.count == 1 || b.count == -1 {
		h := sipHash(b.dataSum)
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
