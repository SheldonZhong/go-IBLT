package iblt

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var tests = []struct {
	dataLen     int
	hashNum     int
	bktNum      uint
	alphaItems  int
	betaItems   int
	sharedItems int
}{
	{4, 4, 80, 20, 30, 20},
	{4, 4, 80, 40, 10, 20},
	{4, 4, 120, 30, 30, 0},
	{4, 4, 1024, 350, 300, 500},
	{4, 4, 1024, 700, 0, 500},
	{4, 4, 1024, 5, 700, 500},
	{4, 4, 1024, 300, 300, 500},
	{16, 4, 1024, 130, 550, 6000},
	{4, 4, 1024, 200, 400, 1000},
}

func runTableTest(t *testing.T, f func(val interface{})) {

}

func TestTable_Insert(t *testing.T) {
	rand.Seed(time.Now().Unix())

	for _, test := range tests {
		b := make([]byte, test.dataLen)
		table := NewTable(test.bktNum, test.dataLen, test.hashNum)
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
		alphaTable := NewTable(test.bktNum, test.dataLen, test.hashNum)
		betaTable := NewTable(test.bktNum, test.dataLen, test.hashNum)
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
		table := NewTable(test.bktNum, test.dataLen, test.hashNum)
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
		table := NewTable(test.bktNum, test.dataLen, test.hashNum)
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
		if !reflect.DeepEqual(rec, cpy) {
			t.Errorf("recoveried IBLT not equal, want %v, get %v", cpy, rec)
		}
	}
}
