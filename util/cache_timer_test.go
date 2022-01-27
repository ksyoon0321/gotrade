package util

import (
	"testing"
	"time"
)

func TestCacheTimer(t *testing.T) {
	c := NewCacheTimer(time.Minute)

	c.Put("11", "22", time.Minute*3)

	v := c.Get("11")

	if v == nil {
		t.Errorf("it should not nil")
	}
	if v.(string) != "22" {
		t.Errorf("it should be 22")
	}

	c2 := NewCacheTimer(time.Second)
	c2.Put("22", 33, time.Second*3)

	v2 := c2.Get("22")
	if v2.(int) != 33 {
		t.Errorf("it should be 33 not %d", v2.(int))
	}

	time.Sleep(time.Second * 4)

	v3 := c2.Get("22")
	if v3 != nil {
		t.Errorf("it should be nil")
	}
}
