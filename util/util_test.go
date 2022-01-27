package util

import (
	"testing"
)

func TestGetAfterPercent(t *testing.T) {
	v := GetAfterPercent(10000, -3)

	if v != 9700 {
		t.Errorf("it should be 9700 not %f", v)
	}
}

func TestBinarySearch(t *testing.T) {
	item := []int{2, 3, 5, 7, 9, 11, 12, 13, 15}

	search := 15
	at := BinarySearch(item, search)
	if item[at] != search {
		t.Errorf("it should be %d not %d, %d", search, item[at], at)
	}

	search = 2
	at = BinarySearch(item, search)
	if item[at] != search {
		t.Errorf("it should be %d not %d, %d", search, item[at], at)
	}

	search = 12
	at = BinarySearch(item, search)
	if item[at] != search {
		t.Errorf("it should be %d not %d, %d", search, item[at], at)
	}

	at = BinarySearch(item, 1)
	if at > 0 {
		t.Errorf("it should be minus not %d", at)
	}

	at = BinarySearch(item, 30)
	if at > 0 {
		t.Errorf("it should be minus2 %d", at)
	}
}

func TestGetBidPrice(t *testing.T) {
	v1 := GetBidPrice(800.0)
	if v1 != 1.0 {
		t.Errorf("it shoulde be 1.0")
	}

	v2 := GetBidPrice(80000.0)
	if v2 != 10.0 {
		t.Errorf("it shoulde be 10.0")
	}
}
