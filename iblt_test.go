package iblt

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var tests = []struct {
	dataLen     int
	hashLen     int
	hashNum     int
	bktNum      uint
	alphaItems  int
	betaItems   int
	sharedItems int
}{
	{4, 1, 4, 80, 20, 30, 20},
	{4, 1, 4, 80, 40, 10, 20},
	{4, 1, 4, 120, 30, 30, 0},
	{4, 1, 4, 1024, 350, 300, 500},
	{4, 1, 4, 1024, 700, 0, 500},
	{4, 1, 4, 1024, 5, 700, 500},
	{4, 1, 4, 1024, 300, 300, 500},
	{16, 1, 4, 1024, 130, 550, 6000},
	{4, 1, 4, 1024, 200, 400, 1000},
}

func TestTable_Insert(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for _, test := range tests {
		b := make([]byte, test.dataLen)
		table := NewTable(test.bktNum, test.dataLen, test.hashLen, test.hashNum)
		for i := 0; i < test.alphaItems; i ++ {
			rand.Read(b)
			if err := table.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}
		diff, err := table.Decode()
		if err != nil {
			t.Errorf("test Decode failed error: %v, case: %v", err, test)
		}
		if diff.AlphaLen() != test.alphaItems {
			t.Errorf("output number of difference mismatch want: %d, get: %d, case: %v", test.alphaItems, diff.AlphaLen(), test)
		}
		if diff.BetaLen() != 0 {
			t.Error("beta diff set is not equal to 0")
		}
	}
}

// IBLT subtract IBLT then decode
func TestTable_Decode(t *testing.T) {
	seed := time.Now().Unix()
	rand.Seed(seed)

	for _, test := range tests {
		alphaTable := NewTable(test.bktNum, test.dataLen, test.hashLen, test.hashNum)
		betaTable := NewTable(test.bktNum, test.dataLen, test.hashLen, test.hashNum)
		b := make([]byte, test.dataLen)
		for i := 0; i < test.alphaItems; i ++ {
			rand.Read(b)
			if err := alphaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}

		for i := 0; i < test.betaItems; i ++ {
			rand.Read(b)
			if err := betaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}

		for i := 0; i < test.sharedItems; i ++ {
			rand.Read(b)
			if err := alphaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
			if err := betaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}

		if err := alphaTable.Subtract(betaTable); err != nil {
			t.Errorf("subtract error: %v", err)
		}

		diff, err := alphaTable.Decode()
		if err != nil {
			t.Errorf("test Decode failed error: %v, case: %v", err, test)
		}

		if diff.AlphaLen() != test.alphaItems {
			t.Errorf("decode diff number mismatched alpha want %d, get %d, case: %v", test.alphaItems, diff.AlphaLen(), test)
		}
		if diff.BetaLen() != test.betaItems {
			t.Errorf("decode diff number mismatched beta want %d, get %d, case :%v", test.betaItems, diff.BetaLen(), test)
		}
	}
}

// construct IBLT and delete one by one and decode
func TestTable_Delete(t *testing.T) {
	seed := time.Now().Unix()
	rand.Seed(seed)

	for _, test := range tests {
		table := NewTable(test.bktNum, test.dataLen, test.hashLen, test.hashNum)
		b := make([]byte, test.dataLen)
		for i := 0; i < test.alphaItems; i ++ {
			rand.Read(b)
			if err := table.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}
		for i := 0; i < test.betaItems; i ++ {
			rand.Read(b)
			if err := table.Delete(b); err != nil {
				t.Errorf("test Delete failed error: %v", err)
			}
		}
		for i := 0; i < test.sharedItems; i ++ {
			rand.Read(b)
			if err := table.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}

			// simulate insert and delete shared items
			if err := table.Delete(b); err != nil {
				t.Errorf("test Delete failed error: %v", err)
			}
		}

		diff, err := table.Decode()
		if err != nil {
			t.Errorf("test Decode failed error: %v, case: %v", err, test)
		}
		if diff.AlphaLen() != test.alphaItems {
			t.Errorf("decode diff number mismatched alpha want %d, get %d, case: %v", test.alphaItems, diff.AlphaLen(), test)
		}
		if diff.BetaLen() != test.betaItems {
			t.Errorf("decode diff number mismatched beta want %d, get %d, case :%v", test.betaItems, diff.BetaLen(), test)
		}
	}
}

func TestTableEncodeDecode(t *testing.T) {
	seed := time.Now().Unix()
	rand.Seed(seed)

	for _, test := range tests {
		table := NewTable(test.bktNum, test.dataLen, test.hashLen, test.hashNum)
		b := make([]byte, test.dataLen)
		for i := 0; i < test.alphaItems; i ++ {
			rand.Read(b)
			if err := table.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}
		for i := 0; i < test.betaItems; i ++ {
			rand.Read(b)
			if err := table.Delete(b); err != nil {
				t.Errorf("test Delete failed error: %v", err)
			}
		}
		cpy := table.Copy()

		enc, err := table.Serialize()
		if err != nil {
			t.Errorf("table serialize error %v", err)
		}

		rec, err := Deserialize(enc)
		if err != nil {
			t.Errorf("recovery from bytes error %v", err)
		}

		if rec.bktNum != cpy.bktNum {
			t.Errorf("recoveried bktNum not equal, want %v, get %v", cpy.bktNum, rec.bktNum)
		}
		if rec.dataLen != cpy.dataLen {
			t.Errorf("recoveried dataLen not equal, want %v, get %v", cpy.dataLen, rec.dataLen)
		}
		if rec.hashLen != cpy.hashLen {
			t.Errorf("recoveried hashLen not equal, want %v, get %v", cpy.hashLen, rec.hashLen)
		}
		if rec.hashNum != cpy.hashNum {
			t.Errorf("recoveried hashNum not equal, want %v, get %v", cpy.hashNum, rec.hashNum)
		}
		for idx, bkt := range rec.buckets {
			cpyBkt := cpy.buckets[idx]
			if bkt == nil && cpyBkt != nil {
				t.Errorf("recoveried bucket is nil at %d, want %v", idx, cpyBkt)
			}
			if bkt != nil && cpyBkt == nil {
				t.Errorf("recoveried bucket is not nil at %d, get %v", idx, bkt)
			}
			if bkt != nil && cpyBkt != nil {
				if bkt.count != cpyBkt.count {
					t.Errorf("recoveried bucket count not equal at %d, want %v, get %v", idx, cpyBkt.count, bkt.count)
				}
				if !bytes.Equal(bkt.dataSum, cpyBkt.dataSum) {
					t.Errorf("recoveried bucket dataSum not equal at %d, want, %v, get %v", idx, cpyBkt.dataSum, bkt.dataSum)
				}
				if !bytes.Equal(bkt.hashSum, cpyBkt.hashSum) {
					t.Errorf("recoveried bucket hashSum not equal at %d, want, %v, get %v", idx, cpyBkt.hashSum, bkt.hashSum)
				}
			}
		}
		if !reflect.DeepEqual(rec, cpy) {
			t.Errorf("recoveried IBLT not equal, want %v, get %v", cpy, rec)
		}
	}
}
