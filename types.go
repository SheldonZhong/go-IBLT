package iblt

import (
	"bytes"
	"fmt"
	"github.com/dchest/siphash"
	"github.com/pkg/errors"
	"github.com/seiflotfy/cuckoofilter"
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

func (b *Bucket) operate(d data, sign bool) {
	b.dataSum.xor(d)
	h := sipHash(d)
	b.hashSum.xor(hash(h))
	if sign {
		b.count++
	} else {
		b.count--
	}
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

type byteSet struct {
	set    [][]byte
	filter *cuckoo.Filter
}

func newByteSet(cap uint) *byteSet {
	return &byteSet{
		set:    make([][]byte, 0),
		filter: cuckoo.NewFilter(cap),
	}
}

func (s byteSet) len() int {
	return len(s.set)
}

// TODO: notify the caller by return an error
func (s *byteSet) insert(b []byte) {
	if !s.test(b) {
		s.filter.Insert(b)
		s.set = append(s.set, b)
	}
}

func (s byteSet) test(b []byte) bool {
	if s.filter.Lookup(b) {
		// possibly in set, false positive
		for _, ele := range s.set {
			if bytes.Equal(b, ele) {
				return true
			}
		}
	}
	return false
}

func (s *byteSet) delete(b []byte) {
	s.filter.Delete(b)
	idx := 0
	for i, ele := range s.set {
		if bytes.Equal(b, ele) {
			idx = i
			break
		}
	}
	s.set = append(s.set[:idx], s.set[idx+1:]...)
}

// each part of symmetric difference
type Diff struct {
	alpha *byteSet
	beta  *byteSet
}

// bktNum as a good estimation for cuckoo filter capacity
func NewDiff(bktNum uint) *Diff {
	return &Diff{
		alpha: newByteSet(bktNum),
		beta:  newByteSet(bktNum),
	}
}

// assume b is pure bucket
func (d *Diff) encode(b *Bucket) error {
	cpy := make([]byte, len(b.dataSum))
	copy(cpy, b.dataSum)
	if b.count == 1 {
		if d.beta.test(cpy) {
			d.beta.delete(cpy)
			return errors.New("repetitive bytes found in beta")
		}
		d.alpha.insert(cpy)
	}
	if b.count == -1 {
		if d.alpha.test(cpy) {
			d.alpha.delete(cpy)
			return errors.New("repetitive bytes found in alpha")
		}
		d.beta.insert(cpy)
	}
	return nil
}

func (d Diff) AlphaItems() int {
	return d.alpha.len()
}

func (d Diff) BetaItems() int {
	return d.beta.len()
}
