package coin

import (
	"testing"

	"github.com/ksyoon0321/gotrade/util"
)

func TestLatest(t *testing.T) {
	q := util.NewCircleQueue(20)

	map1 := make(map[string]interface{})
	map1["a"] = 1.1

	map2 := make(map[string]interface{})
	map2["a"] = 1.3

	map3 := make(map[string]interface{})
	map3["a"] = 1.5

	map4 := make(map[string]interface{})
	map4["a"] = 1.1

	map5 := make(map[string]interface{})
	map5["a"] = 1.2

	map6 := make(map[string]interface{})
	map6["a"] = 1.7

	map7 := make(map[string]interface{})
	map7["a"] = 1.4

	map8 := make(map[string]interface{})
	map8["a"] = 1.6

	map9 := make(map[string]interface{})
	map9["a"] = 1.9

	map10 := make(map[string]interface{})
	map10["a"] = 1.8

	q.Enqueue(map1)
	q.Enqueue(map2)
	q.Enqueue(map3)
	q.Enqueue(map4)
	q.Enqueue(map5)
	q.Enqueue(map6)
	q.Enqueue(map7)
	q.Enqueue(map8)
	q.Enqueue(map9)
	q.Enqueue(map10)

	// low1 := latestHighLow(q, true, 11, "a")
	// if low1 > 0 {
	// 	t.Errorf("it should be 0 not %f", low1)
	// }

	low2 := latestHighLow(q, true, 9, "a")
	if low2 != 1.9 {
		t.Errorf("it should be 0 not %f", low2)
	}

	low3 := latestHighLow(q, false, 9, "a")
	if low3 != 1.1 {
		t.Errorf("it should be 0 not %f", low3)
	}
}
