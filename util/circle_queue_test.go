package util

import (
	"testing"
)

func TestQueue(t *testing.T) {
	count := 2
	q := NewCircleQueue(4)

	q.Enqueue(1)
	q.Enqueue(2)
	/*	q.Enqueue(3)
		q.Enqueue(4)
		q.Enqueue(5)
		q.Enqueue(6)
	*/
	rlt := q.PeekArrayFromTail(count)
	lenrlt := len(rlt)
	if len(rlt) != count {
		t.Errorf(" should be 4 not %d", lenrlt)
	}

	if rlt[0].(int) != 2 {
		t.Errorf(" should be 2 not %d", rlt[0])
	}
	if rlt[1].(int) != 1 {
		t.Errorf(" should be 1 not %d", rlt[1])
	}
}

func TestExportArray(t *testing.T) {
	count := 2
	q := NewCircleQueue(4)

	q.Enqueue(1)
	q.Enqueue(2)
	rlt := q.ExportToArray()
	lenrlt := len(rlt)
	if len(rlt) != count {
		t.Errorf(" should be 4 not %d", lenrlt)
	}

	if rlt[0].(int) != 1 {
		t.Errorf(" should be 1 not %d", rlt[0])
	}
	if rlt[1].(int) != 2 {
		t.Errorf(" should be 2 not %d", rlt[1])
	}
}
