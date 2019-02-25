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

type bucket struct {
	dataSum data
	hashSum hash
	count   int
}

func NewBucket(len int) bucket {
	return bucket{
		dataSum: make(data, len),
		hashSum: hash(0),
		count:   0,
	}
}

func (b *bucket) xor(a bucket) {
	b.dataSum.xor(a.dataSum)
	b.hashSum.xor(a.hashSum)
}

func (b *bucket) add(a bucket) {
	b.xor(a)
	b.count = b.count + a.count
}

func (b *bucket) subtract(a bucket) {
	b.xor(a)
	b.count = b.count - a.count
}

func (b *bucket) put(d data) {
	b.dataSum.xor(d)
	h := sipHash(d)
	b.hashSum.xor(hash(h))
	b.count++
}

func (b bucket) pure() bool {
	if b.count == 1 || b.count == -1 {
		h := sipHash(b.dataSum)
		if b.hashSum == hash(h) {
			return true
		}
	}
	return false
}

func (b bucket) empty() bool {
	return b.count == 0 &&
		b.hashSum.empty() &&
		b.dataSum.empty()
}

func (b bucket) String() string {
	return fmt.Sprintf("bucket: dataSum: %v, hashSum: %v, count: %d",
		b.dataSum, b.hashSum, b.count)
}
