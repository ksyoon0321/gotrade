package tree

import "testing"

func sampleBTree() *BTree {
	b := NewBTree()
	b.Put("11", 1)
	b.Put("06", 3)

	// p = b.PrintTree()
	// t.Errorf("%s", p)

	b.Put("08", 2)
	// p = b.PrintTree()
	// t.Errorf("%s", p)

	b.Put("19", 4)
	b.Put("04", 5)

	b.Put("10", 6)

	b.Put("05", 77)
	b.Put("17", 44)
	b.Put("43", 55)
	b.Put("49", 551)
	b.Put("32", 552)
	b.Put("31", 553)
	b.Put("30", 554)
	b.Put("33", 555)
	b.Put("45", 556)
	b.Put("45", 557)
	b.Put("12", 128)
	b.Put("09", 129)

	return b
}
func TestBTree(t *testing.T) {

	b := sampleBTree()
	v1 := b.Get("12")
	if v1 == nil {
		t.Errorf("it should not nil v1")
	}

	if v1.(int) != 128 {
		t.Errorf("it should be 128 not %d", v1.(int))
	}

	//	p = b.PrintTree()
	//	t.Errorf("%s", p)

}

func TestBTreePrint(t *testing.T) {
	p := ""
	b := sampleBTree()
	p = b.PrintTree()
	t.Errorf("%s", p)

}
