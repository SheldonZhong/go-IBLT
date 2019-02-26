package iblt

import (
	"errors"
	"github.com/dchest/siphash"
	"github.com/willf/bitset"
)

type Table struct {
	dataLen int
	hashNum int
	bktNum  uint
	buckets []*Bucket
	bitsSet *bitset.BitSet
}

// Specify number of buckets, data field length (in byte), number of hash functions
func NewTable(buckets uint, dataLen int, hashNum int, ) *Table {
	return &Table{
		dataLen: dataLen,
		hashNum: hashNum,
		bktNum:  buckets,
		// TODO: save space by embedding index in bucket
		buckets: make([]*Bucket, buckets),
		bitsSet: bitset.New(buckets),
	}
}

func (t *Table) Insert(d []byte) error {
	if len(d) != t.dataLen {
		return errors.New("mismatched data length")
	}

	t.bitsSet.ClearAll()
	tries := 0
	for i := 0; i < t.hashNum; {
		// assume we can always find different keys
		// as this is in high probability
		h := siphash.Hash(key0, uint64(key1+tries), d)
		tries++
		// TODO: modulate produces imbalanced uniform distribution
		idx := uint(h) % t.bktNum
		if !t.bitsSet.Test(idx) {
			t.put(idx, d)
			t.bitsSet.Set(idx)
			i++
		}
	}
	return nil
}

func (t *Table) put(idx uint, d []byte) {
	if t.buckets[idx] == nil {
		t.buckets[idx] = NewBucket(t.dataLen)
	}
	t.buckets[idx].put(d)
}
