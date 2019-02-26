package iblt

type Table struct {
	hashLen int
	buckets []Bucket
}

func NewTable(buckets int, hashLen int) *Table {
	return &Table{
		hashLen: hashLen,
		buckets: make([]Bucket, buckets),
	}
}

func (t *Table) add() {

}
