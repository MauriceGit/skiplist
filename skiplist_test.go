package skiplist

import (
	"fmt"
	"math"
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

type FloatElement float64

func (e FloatElement) ExtractKey() float64 {
	return float64(e)
}
func (e FloatElement) String() string {
	return fmt.Sprintf("%.3f", e)
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
	eps := 0.00000001
	list := NewEps(eps)
	maxNumber := 1000.0

	for i := 0; i < maxN; i++ {
		list.Insert(FloatElement(rand.Float64() * maxNumber))
	}

	first := float64(list.GetSmallestNode().GetValue().(FloatElement))

	for i := 0; i < maxN; i++ {
		f := rand.Float64() * maxNumber
		if v, ok := list.FindGreaterOrEqual(FloatElement(f)); ok {
			// if f is v should be bigger than the element before
			lastV := float64(list.Prev(v).GetValue().(FloatElement))
			thisV := float64(v.GetValue().(FloatElement))
			isFirst := math.Abs(first-thisV) <= eps
			if !isFirst && lastV >= f {
				fmt.Printf("PrevV: %.8f\n    f: %.8f\n\n", lastV, f)
				t.Fail()
			}
			// v should be bigger or equal to f
			// If we compare directly, we get an equal key with a difference on the 10th decimal point, which fails.
			if f-thisV > eps {
				fmt.Printf("f: %.8f\nv: %.8f\n\n", f, thisV)
				t.Fail()
			}
		} else {
			lastV := float64(list.GetLargestNode().GetValue().(FloatElement))
			// It is OK, to fail, as long as f is bigger than the last element.
			if f <= lastV {
				fmt.Printf("lastV: %.8f\n    f: %.8f\n\n", lastV, f)
				t.Fail()
			}
		}
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
		// Next.Prev must always point to itself!
		if list.Prev(list.Next(node)) != node {
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
		// Next.Prev must always point to itself!
		if list.Next(list.Prev(node)) != node {
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
