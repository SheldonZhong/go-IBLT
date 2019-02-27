package iblt

import (
	"math/rand"
	"testing"
)

func TestTable_Insert(t *testing.T) {
	tests := []struct {
		dataLen int
		hashNum int
		bktNum  uint
		items   int
	}{
		{4, 4, 80, 59},
		{4, 4, 80, 60},
		{4, 4, 120, 80},
		{4, 4, 1024, 700},
	}

	for _, test := range tests {
		b := make([]byte, test.dataLen)
		table := NewTable(test.bktNum, test.dataLen, test.hashNum)
		for i := 0; i < test.items; i ++ {
			rand.Read(b)
			if err := table.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
		}
		diff, err := table.Decode()
		if err != nil {
			t.Errorf("test Decode failed error: %v", err)
		}
		if len(diff.Alpha) != test.items {
			t.Errorf("output number of difference mismatch want: %d, get: %d", test.items, len(diff.Alpha))
		}
	}
}

func TestTable_Decode(t *testing.T) {
	tests := []struct {
		dataLen     int
		hashNum     int
		bktNum      uint
		alphaItems  int
		betaItems   int
		sharedItems int
	}{
		{4, 4, 80, 20, 20, 100},
		{4, 4, 80, 20, 20, 200},
		{4, 4, 120, 30, 30, 300},
		{4, 4, 1024, 50, 50, 500},
	}

	for _, test := range tests {
		alphaBuff := make([][]byte, test.alphaItems)
		betaBuff := make([][]byte, test.betaItems)

		alphaTable := NewTable(test.bktNum, test.dataLen, test.hashNum)
		betaTable := NewTable(test.bktNum, test.dataLen, test.hashNum)
		for i := 0; i < test.alphaItems; i ++ {
			b := make([]byte, test.dataLen)
			rand.Read(b)
			if err := alphaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
			alphaBuff = append(alphaBuff, b)
		}

		for i := 0; i < test.betaItems; i ++ {
			b := make([]byte, test.dataLen)
			rand.Read(b)
			if err := betaTable.Insert(b); err != nil {
				t.Errorf("test Insert failed error: %v", err)
			}
			betaBuff = append(betaBuff, b)
		}

		for i := 0; i < test.sharedItems; i ++ {
			b := make([]byte, test.dataLen)
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
			t.Errorf("test Decode failed error: %v", err)
		}
		if len(diff.Alpha) != test.alphaItems {
			t.Errorf("decode diff number mismatched alpha want %d, get %d", test.alphaItems, len(diff.Alpha))
		}
		if len(diff.Beta) != test.betaItems {
			t.Errorf("decode diff number mismatched beta want %d, get %d", test.betaItems, len(diff.Beta))
		}
	}
}
