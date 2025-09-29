package arbitrage

import (
	"testing"

	"github.com/0xarash/arbitrage/src/decimal"
)

func TestPriority(t *testing.T) {
	pq := NewPriorityQueue[string]()

	n1, _ := decimal.New(5, 5)
	n2, _ := decimal.New(1, 5)
	n3, _ := decimal.New(4, 5)
	n4, _ := decimal.New(9, 5)
	n5, _ := decimal.New(10, 5)
	n6, _ := decimal.New(3, 5)
	n7, _ := decimal.New(6, 5)
	n8, _ := decimal.New(20, 5)
	n9, _ := decimal.New(100, 5)
	n10, _ := decimal.New(2, 5)

	pq.Push("5", n1)
	pq.Push("1", n2)
	pq.Push("4", n3)
	pq.Push("9", n4)
	pq.Push("10", n5)
	pq.Push("3", n6)
	pq.Push("6", n7)
	pq.Push("20", n8)
	pq.Push("100", n9)
	pq.Push("2", n10)

	v1 := pq.Pop()
	expect := "100"
	if v1 != expect {
		t.Errorf("Got %v Expect %v", v1, expect)
	}

	v2 := pq.Pop()
	expect = "20"
	if v2 != expect {
		t.Errorf("Got %v Expect %v", v2, expect)
	}

	v3 := pq.Pop()
	expect = "10"
	if v3 != expect {
		t.Errorf("Got %v Expect %v", v3, expect)
	}

	v4 := pq.Pop()
	expect = "9"
	if v4 != expect {
		t.Errorf("Got %v Expect %v", v4, expect)
	}

	v5 := pq.Pop()
	expect = "6"
	if v5 != expect {
		t.Errorf("Got %v Expect %v", v5, expect)
	}
}
