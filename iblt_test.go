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
		{4, 4, 80, 10},
		{4, 4, 80, 13},
		{4, 4, 120, 25},
		{4, 4, 1024, 50},
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
