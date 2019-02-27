package iblt

import (
	"errors"
	"github.com/dchest/siphash"
	"github.com/golang-collections/collections/queue"
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
	err := t.index(d)
	if err != nil {
		return err
	}

	for i, e := t.bitsSet.NextSet(0); e; i, e = t.bitsSet.NextSet(i + 1) {
		t.put(i, d)
	}
	return nil
}

func (t *Table) Delete(d []byte) error {
	cpy := make([]byte, len(d))
	copy(cpy, d)
	err := t.index(cpy)
	if err != nil {
		return err
	}

	for i, e := t.bitsSet.NextSet(0); e; i, e = t.bitsSet.NextSet(i + 1) {
		t.take(i, cpy)
	}
	return nil
}

func (t *Table) index(d []byte) error {
	if len(d) != t.dataLen {
		return errors.New("insert byte length mismatches base data length")
	}

	if t.bitsSet == nil {
		t.bitsSet = bitset.New(t.bktNum)
	}

	t.bitsSet.ClearAll()
	tries := 0
	for i := 0; i < t.hashNum; {
		// assume we can always find different keys
		// as this is in high probability
		h := siphash.Hash(key0, uint64(key1+tries), d)
		tries++
		// TODO: modulo produces imbalanced uniform distribution
		idx := uint(h) % t.bktNum
		if !t.bitsSet.Test(idx) {
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

// Decode is self-destructive
func (t *Table) Decode() (*Diff, error) {
	pure := queue.New()
	pureMask := bitset.New(t.bitsSet.Len())
	err := t.enqueuePure(pureMask, pure)
	if err != nil {
		return nil, err
	}
	// ensure we have at least one pure bucket in the IBLT
	// this is necessary condition for decoding an IBLT
	if pure.Len() == 0 {
		return nil, errors.New("no pure buckets in table")
	}

	diff := NewDiff()
	bkt := NewBucket(t.dataLen)
	for pure.Len() > 0 {
		// clean out pure queue, delete all pure buckets and output the stored data
		// it will create more pure buckets to decode in the next cycle
		for pure.Len() > 0 {
			bkt = pure.Dequeue().(*Bucket)
			diff.encode(bkt)
			// Delete would possibly delete pure buckets that stored in queue before
			err = t.Delete(bkt.dataSum)
			if err != nil {
				return diff, err
			}
		}
		// now pure queue should be empty, enqueue more pure cell
		err = t.enqueuePure(pureMask, pure)
		if err != nil {
			return diff, err
		}
		// no more bucket is pure either
		// 1) we have successfully decoded all the possible buckets and all the buckets should be empty
		// 2) we have hash collision for more than two items
	}
	// check if every bucket is empty
	for i := range t.buckets {
		if t.buckets[i] != nil && !t.buckets[i].empty() {
			return diff, errors.New("dirty entries remained")
		}
	}

	return diff, nil
}

func (t *Table) enqueuePure(pureMask *bitset.BitSet, pure *queue.Queue) error {
	for i := range t.buckets {
		// skip the same pure bucket at difference indexes, enqueue the first one
		if !pureMask.Test(uint(i)) && t.buckets[i] != nil && t.buckets[i].pure() {
			err := t.index(t.buckets[i].dataSum)
			if err != nil {
				return err
			}
			pureMask.InPlaceUnion(t.bitsSet)
			pure.Enqueue(t.buckets[i])
		}
	}
	return nil
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

func (t *Table) take(idx uint, d []byte) {
	if t.buckets[idx] == nil {
		t.buckets[idx] = NewBucket(t.dataLen)
	}
	t.buckets[idx].take(d)
}
