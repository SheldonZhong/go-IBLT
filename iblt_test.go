package iblt

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTable_Insert(t *testing.T) {
	rand.Seed(time.Now().Unix())
	tests := []struct {
		dataLen int
		hashNum int
		bktNum  uint
		items   int
	}{
		{4, 4, 40, 20},
		{4, 1, 40, 2},
		{4, 2, 40, 5},
		{4, 3, 40, 10},
		{4, 5, 40, 10},
		{4, 6, 40, 10},
		{4, 4, 80, 50},
		{6, 4, 120, 70},
		{8, 4, 1024, 700},
		{20, 4, 2000, 1300},
		{32, 4, 4000, 300},
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
		if diff.AlphaItems() != test.items {
			t.Errorf("output number of difference mismatch want: %d, get: %d, case: %v", test.items, diff.AlphaItems(), test)
		}
	}
}

func TestTable_Decode(t *testing.T) {
	seed := time.Now().Unix()
	//seed := int64(1551433058)
	rand.Seed(seed)
	fmt.Println(seed)
	tests := []struct {
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
		{4, 4, 1024, 0, 700, 500},
		{4, 4, 1024, 300, 300, 500},
		{16, 4, 1024, 130, 550, 6000},
		{4, 4, 1024, 200, 400, 1000},
	}

	for _, test := range tests {
		alphaBuff := make([][]byte, 0)
		betaBuff := make([][]byte, 0)

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
			t.Errorf("test Decode failed error: %v, case: %v", err, test)
		}
		//debugBucket(t, alphaTable)

		if diff.AlphaItems() != test.alphaItems {
			bytesCompare(alphaBuff, diff.alpha.set)
			t.Errorf("decode diff number mismatched alpha want %d, get %d, case: %v", test.alphaItems, diff.AlphaItems(), test)
		}
		if diff.BetaItems() != test.betaItems {
			bytesCompare(betaBuff, diff.beta.set)
			t.Errorf("decode diff number mismatched beta want %d, get %d, case :%v", test.betaItems, diff.BetaItems(), test)
		}
		fmt.Println("------------test case ends------------")
	}
}

// iterate over beta, for each element print out those does not exist in alpha
func bytesCompare(alpha [][]byte, beta [][]byte) bool {
	allFound := true
	fmt.Print("extra ")
	for _, b := range beta {
		found := false
		for _, a := range alpha {
			if bytes.Compare(a, b) == 0 {
				found = true
				break
			}
		}

		if !found {
			fmt.Println(b)
			allFound = false
		}
	}
	return allFound
}

func TestHash(t *testing.T) {
	table := NewTable(1024, 4, 4)
	if err := table.index([]byte{131, 250, 218, 247}); err != nil {
		t.Errorf("error index")
	}

	print := ""
	for i, e := table.bitsSet.NextSet(0); e; i, e = table.bitsSet.NextSet(i + 1) {
		print += fmt.Sprintf("%d ", i)
	}

	fmt.Println(print)
}

func debugBucket(t *testing.T, table *Table) {
	for _, bkt := range table.buckets {
		fmt.Println(bkt)
	}
}
