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
		return errors.New("insert byte length mismatches base data length")
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

// Modify callee
func (t *Table) Subtract(a *Table) error {
	err := t.check(a)
	if err != nil {
		return err
	}

	for i := range t.buckets {
		if t.buckets[i] != nil && a.buckets[i] != nil {
			t.buckets[i].subtract(a.buckets[i])
		}
		if t.buckets[i] == nil && a.buckets[i] != nil {
			t.buckets[i] = a.buckets[i].copy()
			t.buckets[i].count = -t.buckets[i].count
		}
	}

	return err
}

func (t Table) check(a *Table) error {
	if t.bktNum != a.bktNum {
		return errors.New("subtract table mismatches bucket number")
	}

	if t.dataLen != a.dataLen {
		return errors.New("subtract table mismatches data length")
	}

	if t.hashNum != a.hashNum {
		return errors.New("subtract table mismatches number of hash functions")
	}

	if len(t.buckets) != len(a.buckets) {
		return errors.New("illegally appended buckets")
	}

	return nil
}

func (t *Table) put(idx uint, d []byte) {
	if t.buckets[idx] == nil {
		t.buckets[idx] = NewBucket(t.dataLen)
	}
	t.buckets[idx].put(d)
}
