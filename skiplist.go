// MIT License
//
// Copyright (c) 2018 Maurice Tollmien (maurice.tollmien@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// tree23 is an implementation for a balanced 2-3-tree.
// It distinguishes itself from other implementations of 2-3-trees by having a few more
// functions defined for finding elements close to a key (similar to possible insert positions in the tree)
// for floating point keys and by having a native function to retreive the next or previous leaf SkipListElement
// in the tree without knowing its key or position in the tree that work in O(1) for every leaf!
// The last SkipListElement links to the first and the first back to the last SkipListElement.
// The tree has its own memory manager to avoid frequent allocations for single nodes that are created or removed.
package skiplist

import (
    "fmt"
    "math/bits"
    "math/rand"
    "math"
    "time"
)

const (
    MAX_LEVEL = 16
)

type ListElement interface {
    ExtractValue() float64
    String() string
}

type SkipListPointer struct {
    prev *SkipListElement
    next *SkipListElement
}

type Backtrack struct {
    node    *SkipListElement
    level   int
}

type SkipListElement struct {
    array       [MAX_LEVEL]SkipListPointer
    level       int
    key         float64
    value       ListElement
}

type SkipList struct {
    startLevels         [MAX_LEVEL]*SkipListElement
    endLevels           [MAX_LEVEL]*SkipListElement
    backtrack           [MAX_LEVEL]Backtrack
    lastBacktrackCount  int
    maxNewLevel         int
    maxLevel            int
    elementCount        int
    elementSum          float64
    eps                 float64
}

// Package initialization
func init() {
    seed := time.Now().UTC().UnixNano()
    //seed = 1530076445104807822
    fmt.Printf("seed: %v\n", seed)
    rand.Seed(seed)
}

func generateLevel(maxLevel int) int {
    // First we apply some mask which makes sure that we don't get a level
    // above our desired level. Then we find the first set bit.
    var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) -1)
    zeroes := bits.TrailingZeros64(x)
    if zeroes <= maxLevel {
        return zeroes
    }
    return maxLevel-1
}

func New(eps float64) SkipList {
    return SkipList{
        startLevels:        [MAX_LEVEL]*SkipListElement{},
        endLevels:          [MAX_LEVEL]*SkipListElement{},
        backtrack:          [MAX_LEVEL]Backtrack{},
        lastBacktrackCount: 0,
        maxNewLevel:        MAX_LEVEL,
        maxLevel:           0,
        elementCount:       0,
        elementSum:         0.0,
        eps:                eps,
    }
}

func (t *SkipList) isEmpty() bool {
    return t.startLevels[0] == nil
}

// returns: found element, backtracking list: Includes the elements from the entry point down to the element (or possible insertion position)!, ok, if an element was found
func (t *SkipList) findExtended(key float64, findGreaterOrEqual bool, createBackTrack bool) (foundElem *SkipListElement, ok bool) {


    foundElem = nil
    ok = false
    increasingSearch := true

    if t.isEmpty() {
        return
    }

    // Find out, if it makes more sense, to search from the left or the right side!
    // Lets just test this feature first, when there is no backtrack created. So just for find itself.
    // I decided, that the effect of starting from the right side is nearly non-observable. It only really shows, when we look for one of the
    // very last elements. So OK for finding but not worth it to use in inserts as it would complicate things a lot!
    avg := t.elementSum/float64(t.elementCount)
    if !createBackTrack && key > avg {
        increasingSearch = false
    }

    if createBackTrack {
        t.lastBacktrackCount = 0
    }

    index := 0
    var currentNode *SkipListElement = nil

    // Find good entry point so we don't accidently skip half the list.
    for i := t.maxLevel; i >= 0; i-- {
        if increasingSearch {
            if t.startLevels[i] != nil && t.startLevels[i].key <= key {
                index = i
                break
            }
        } else {
            if t.endLevels[i] != nil && t.endLevels[i].key >= key {
                index = i
                break
            }
        }
    }
    if increasingSearch {
        currentNode = t.startLevels[index]
    } else {
        currentNode = t.endLevels[index]
    }

    currCompare := 1
    if currentNode.key < key {
        currCompare = -1
    } else if math.Abs(currentNode.key - key) <= t.eps {
        currCompare = 0
    }

    nextCompare := 0

    for {
        if currCompare == 0 {
            foundElem = currentNode
            ok = true
            return
        }

        nextNode := currentNode.array[index].next
        if !increasingSearch {
            nextNode = currentNode.array[index].prev
        }

        if nextNode != nil {
            nextCompare = 1
            if nextNode.key < key {
                nextCompare = -1
            } else if math.Abs(nextNode.key - key) <= t.eps {
                nextCompare = 0
            }
            currCompare = nextCompare
        }

        // Which direction are we continuing next time?
        if nextNode != nil && (increasingSearch && nextCompare <= 0 || !increasingSearch && nextCompare >= 0) {
            // Go right
            currentNode = nextNode
        } else {
            if createBackTrack {
                t.backtrack[t.lastBacktrackCount].node = currentNode
                t.backtrack[t.lastBacktrackCount].level = index
                t.lastBacktrackCount++
            }
            if index > 0 {
                // Go down
                index--
            } else {
                // Element is not found and we reached the bottom.
                if findGreaterOrEqual {
                    foundElem = nextNode
                    ok = nextNode != nil
                }
                return

            }
        }
    }

    return
}

func (t *SkipList) Find(e ListElement) (*SkipListElement, bool) {
    return t.findExtended(e.ExtractValue(), false, false)
}

func (t *SkipList) FindGreaterOrEqual(e ListElement) (*SkipListElement, bool) {
    return t.findExtended(e.ExtractValue(), true, false)
}

func (t *SkipList) Delete(e ListElement) {

    if elem,ok := t.Find(e); ok {
        for i := elem.level; i >= 0; i-- {
            prev := elem.array[i].prev
            next := elem.array[i].next

            if prev != nil {
                prev.array[i].next = next
            }
            if next != nil {
                next.array[i].prev = prev
            }

            if t.startLevels[i] == elem {
                t.startLevels[i] = next
                if next == nil {
                    // reduce the maximum entry position!
                    t.maxLevel = i-1
                }
            }
            if t.endLevels[i] == elem {
                t.endLevels[i] = prev
            }
        }
        t.elementCount--
        t.elementSum -= elem.key
    }
}

func (t *SkipList) Insert(e ListElement) {

    level := generateLevel(t.maxNewLevel)
    elem  := &SkipListElement {
                array: [MAX_LEVEL]SkipListPointer{},
                level: level,
                key:   e.ExtractValue(),
                value: e,
            }

    t.elementCount++
    t.elementSum += elem.key

    newFirst := true
    newLast := true
    if !t.isEmpty() {
        newFirst = elem.key < t.startLevels[0].key
        newLast  = elem.key > t.endLevels[0].key
    }

    normallyInserted := false
    // Insertion using Find()
    if !newFirst && !newLast {

        normallyInserted = true

        // Search for e down to level 1. It will not find anything, but will return a backtrack for insertion.
        // We only care about the backtracking anyway.
        t.findExtended(elem.key, true, true)

        btCount := t.lastBacktrackCount

        i := btCount-1
        for i = btCount-1; i >= 0; i-- {

            bt := t.backtrack[i]

            if bt.level > elem.level {
                break
            }

            oldNext := bt.node.array[bt.level].next
            if oldNext != nil {
                oldNext.array[bt.level].prev = elem
            }
            elem.array[bt.level].next = oldNext
            elem.array[bt.level].prev = bt.node
            bt.node.array[bt.level].next = elem
        }
    }

    if level > t.maxLevel {
        t.maxLevel = level
    }

    // Where we have a left-most position that needs to be referenced!
    for  i := level; i >= 0; i-- {

        didSomething := false

        if newFirst || normallyInserted  {
            if elem.array[i].prev == nil {
                if t.startLevels[i] != nil {
                    t.startLevels[i].array[i].prev = elem
                }
                elem.array[i].next = t.startLevels[i]
                t.startLevels[i] = elem
            }

            // link the endLevels to this element!
            if elem.array[i].next == nil {
                t.endLevels[i] = elem
            }

            didSomething = true
        }

        if newLast {
            // Places the element after the very last element on this level!
            // This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
            if !newFirst {
                if t.endLevels[i] != nil {
                    t.endLevels[i].array[i].next = elem
                }
                elem.array[i].prev = t.endLevels[i]
                t.endLevels[i] = elem
            }

            // Link the startLevels to this element!
            if elem.array[i].prev == nil {
                t.startLevels[i] = elem
            }

            didSomething = true
        }

        if !didSomething {
            break
        }
    }

}

func (t *SkipList) PrettyPrint() {

    fmt.Printf(" --> ")
    for i,l := range t.startLevels {
        next := "---"
        if l != nil {
            next = l.value.String()
        }
        fmt.Printf("[%v]    ", next)
        if i < len(t.startLevels)-1 {
            fmt.Printf(" --> ")
        }
    }
    fmt.Println("")

    node := t.startLevels[0]
    for node != nil {
        fmt.Printf("%v: ", node.value)
        for i := 0; i <= node.level; i++ {
            l := node.array[i]

            prev := "---"
            if l.prev != nil {
                prev = l.prev.value.String()
            }
            next := "---"
            if l.next != nil {
                next = l.next.value.String()
            }

            fmt.Printf("[%v|%v]", prev, next)
            if i < node.level {
                fmt.Printf(" --> ")
            }

        }
        fmt.Printf("\n")
        node = node.array[0].next
    }

     fmt.Printf(" --> ")
    for i,l := range t.endLevels {
        next := "---"
        if l != nil {
            next = l.value.String()
        }
        fmt.Printf("[%v]    ", next)
        if i < len(t.endLevels)-1 {
            fmt.Printf(" --> ")
        }
    }
    fmt.Println("")

    fmt.Printf("\n")


}
