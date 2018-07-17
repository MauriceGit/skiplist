package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	//"github.com/pkg/profile"
)

const (
	maxN = 1000000
)

type Element int

func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

type ComplexElement struct {
	E int
	S string
}

func (e ComplexElement) ExtractKey() float64 {
	return float64(e.E)
}
func (e ComplexElement) String() string {
	return fmt.Sprintf("%03d", e.E)
}

// timeTrack will print out the number of nanoseconds since the start time divided by n
// Useful for printing out how long each iteration took in a benchmark
func timeTrack(start time.Time, n int, name string) {
	loopNS := time.Since(start).Nanoseconds() / int64(n)
	fmt.Printf("%s: %d\n", name, loopNS)
}

func TestInsertAndFind(t *testing.T) {
	list := New()
	// Test at the beginning of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(Element(maxN - i))
	}
	for i := 0; i < maxN; i++ {
		if _, ok := list.Find(Element(maxN - i)); !ok {
			t.Fail()
		}
	}

	list = New()
	// Test at the end of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}
	for i := 0; i < maxN; i++ {
		if _, ok := list.Find(Element(i)); !ok {
			t.Fail()
		}
	}

	list = New()
	// Test at random positions in the list.
	rList := rand.Perm(maxN)
	for _, e := range rList {
		list.Insert(Element(e))
	}
	for _, e := range rList {
		if _, ok := list.Find(Element(e)); !ok {
			t.Fail()
		}
	}

}

func TestDelete(t *testing.T) {
	list := New()
	// Delete elements at the beginning of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}
	for i := 0; i < maxN; i++ {
		list.Delete(Element(i))
	}
	if !list.IsEmpty() {
		t.Fail()
	}

	list = New()
	// Delete elements at the end of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}
	for i := 0; i < maxN; i++ {
		list.Delete(Element(maxN - i - 1))
	}
	if !list.IsEmpty() {
		t.Fail()
	}

	list = New()
	// Delete elements at random positions in the list.
	rList := rand.Perm(maxN)
	for _, e := range rList {
		list.Insert(Element(e))
	}
	for _, e := range rList {
		list.Delete(Element(e))
	}
	if !list.IsEmpty() {
		t.Fail()
	}
}

func TestFindGreaterOrEqual(t *testing.T) {
	list := New()

	for i := 0; i < maxN; i++ {
		if i != 45 &&
			i != 46 &&
			i != 47 &&
			i != 48 &&
			i != 6006 &&
			i != 6007 &&
			i != 6001 &&
			i != 6003 {
			list.Insert(Element(i))
		}
	}

	if e, ok := list.FindGreaterOrEqual(Element(44)); ok {
		if e.value.(Element) != 44 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if e, ok := list.FindGreaterOrEqual(Element(45)); ok {
		if e.value.(Element) != 49 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if e, ok := list.FindGreaterOrEqual(Element(47)); ok {
		if e.value.(Element) != 49 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if e, ok := list.FindGreaterOrEqual(Element(6006)); ok {
		if e.value.(Element) != 6008 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if e, ok := list.FindGreaterOrEqual(Element(6001)); ok {
		if e.value.(Element) != 6002 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

	if e, ok := list.FindGreaterOrEqual(Element(6002)); ok {
		if e.value.(Element) != 6002 {
			t.Fail()
		}
	} else {
		t.Fail()
	}

}

func TestPrev(t *testing.T) {
	list := New()

	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}

	smallest := list.GetSmallestNode()
	largest := list.GetLargestNode()

	lastNode := largest
	node := lastNode
	for node != smallest {
		node = list.Prev(node)
		// Must always be incrementing here!
		if node.value.(Element) >= lastNode.value.(Element) {
			t.Fail()
		}
		lastNode = node
	}
}

func TestNext(t *testing.T) {
	list := New()

	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}

	smallest := list.GetSmallestNode()
	largest := list.GetLargestNode()

	lastNode := smallest
	node := lastNode
	for node != largest {
		node = list.Next(node)
		// Must always be incrementing here!
		if node.value.(Element) <= lastNode.value.(Element) {
			t.Fail()
		}
		lastNode = node
	}
}

func TestChangeValue(t *testing.T) {
	list := New()

	for i := 0; i < maxN; i++ {
		list.Insert(ComplexElement{i, "value"})
	}

	for i := 0; i < maxN; i++ {
		// The key only looks at the int so the string doesn't matter here!
		f1, ok := list.Find(ComplexElement{i, ""})
		if !ok {
			t.Fail()
		}
		ok = list.ChangeValue(f1, ComplexElement{i, "different value"})
		if !ok {
			t.Fail()
		}
		f2, ok := list.Find(ComplexElement{i, ""})
		if !ok {
			t.Fail()
		}
		if f2.GetValue().(ComplexElement).S != "different value" {
			t.Fail()
		}
	}

}
