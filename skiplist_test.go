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
	var list SkipList

	var listPointer *SkipList
	listPointer.Insert(Element(0))
	if _, ok := listPointer.Find(Element(0)); ok {
		t.Fail()
	}

	list = New()

	if _, ok := list.Find(Element(0)); ok {
		t.Fail()
	}
	if !list.IsEmpty() {
		t.Fail()
	}

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

	var list SkipList

	// Delete on empty list
	list.Delete(Element(0))

	list = New()

	list.Delete(Element(0))
	if !list.IsEmpty() {
		t.Fail()
	}

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
	maxNumber := 1000.0

	var list SkipList
	var listPointer *SkipList

	// Test on empty list.
	if _, ok := listPointer.FindGreaterOrEqual(FloatElement(0)); ok {
		t.Fail()
	}

	list = NewEps(eps)

	for i := 0; i < maxN; i++ {
		list.Insert(FloatElement(rand.Float64() * maxNumber))
	}

	first := float64(list.GetSmallestNode().GetValue().(FloatElement))

	// Find the very first element. This is a special case in the implementation that needs testing!
	if v, ok := list.FindGreaterOrEqual(FloatElement(first - 2.0*eps)); ok {
		// We found an element different to the first one!
		if math.Abs(float64(v.GetValue().(FloatElement))-first) > eps {
			t.Fail()
		}
	} else {
		// No element found.
		t.Fail()
	}

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

	if list.Prev(smallest) != largest {
		t.Fail()
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

	if list.Next(largest) != smallest {
		t.Fail()
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
		if ok = list.ChangeValue(f2, ComplexElement{i + 5, "different key"}); ok {
			t.Fail()
		}
	}
}

func TestGetNodeCount(t *testing.T) {
	list := New()

	for i := 0; i < maxN; i++ {
		list.Insert(Element(i))
	}

	if list.GetNodeCount() != maxN {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	list := NewSeed(1531889620180049576)

	for i := 0; i < 20; i++ {
		list.Insert(Element(i))
	}

	testString := ` --> [000]     -> [002] -> [009] -> [010]
000: [---|001]
001: [000|002]
002: [001|003] -> [004]
003: [002|004]
004: [003|005] -> [005]
005: [004|006] -> [009]
006: [005|007]
007: [006|008]
008: [007|009]
009: [008|010] -> [010] -> [010]
010: [009|011] -> [012] -> [---] -> [---]
011: [010|012]
012: [011|013] -> [013]
013: [012|014] -> [---]
014: [013|015]
015: [014|016]
016: [015|017]
017: [016|018]
018: [017|019]
019: [018|---]
 --> [019]     -> [013] -> [010] -> [010]
`

	if list.String() != testString {
		t.Fail()
	}
}

func TestInfiniteLoop(t *testing.T) {
	list := New()
	list.Insert(Element(1))

	if _, ok := list.Find(Element(2)); ok {
		t.Fail()
	}

	if _, ok := list.FindGreaterOrEqual(Element(2)); ok {
		t.Fail()
	}
}
