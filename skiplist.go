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
    "time"
)


type ListElement interface {
    Compare(e ListElement) int
    String() string
}

type SkipListPointer struct {
    prev *SkipListElement
    next *SkipListElement
}

type SkipListElement struct {
    array       []SkipListPointer
    level       int
    value       ListElement
}

type SkipList struct {
    levels              [25]*SkipListElement
    maxNewLevel         int
    maxLevel       int
}

// Package initialization
func init() {
    seed := time.Now().UTC().UnixNano()
    seed = 1529743902965759434
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

func New() SkipList {
    return SkipList{[25]*SkipListElement{}, 25, 0}
}

func (t *SkipList) isEmpty() bool {
    return t.levels[0] == nil
}

func (localRoot *SkipListElement) insertRec(e *SkipListElement, height, level int) {

    // Next one is not overshot -- We go right!
    next := localRoot.array[height].next
    if next != nil && next.value.Compare(e.value) < 0 {
        next.insertRec(e, height, level)
        return
    }

    oldNext := localRoot.array[height].next

    // Our level is now the same as height. So we have to squeeze our new SkipListElement in between.
    if level >= height && (oldNext == nil || e.value.Compare(oldNext.value) < 0) {

        if oldNext != nil {
            oldNext.array[height].prev = e
        }
        e.array[height].next = oldNext
        e.array[height].prev = localRoot
        localRoot.array[height].next = e
    }

    if height > 0 {
        localRoot.insertRec(e, height-1, level)
    }
}

func (t *SkipList) Insert(e ListElement) {

    level := generateLevel(t.maxNewLevel)
    elem  := &SkipListElement{make([]SkipListPointer, level+1, level+1), level, e}


    newFirst := true
    if !t.isEmpty() {
        newFirst = t.levels[0].value.Compare(e) > 0
    }
    entryIndex := t.maxLevel
    // Find good entry point so we don't accidently skip half the list.
    if !newFirst {
        for i := t.maxLevel; i >= 0; i-- {
            if t.levels[i] != nil && t.levels[i].value.Compare(e) < 0 {
                entryIndex = i
                break
            }
        }

        if !t.isEmpty() {
            t.levels[entryIndex].insertRec(elem, entryIndex, level)
        }
    }

    if level > t.maxLevel {
        t.maxLevel = level
    }

    // Where we have a left-most position that needs to be referenced!
    for  i := level; i >= 0; i-- {
        if newFirst || elem.array[i].prev == nil {
            if t.levels[i] != nil {
                t.levels[i].array[i].prev = elem
            }
            elem.array[i].next = t.levels[i]
            t.levels[i] = elem

        } else {
            break
        }
    }
}

func (t *SkipList) Find(e ListElement) (*SkipListElement, bool) {
    if t.isEmpty() {
        return nil, false
    }

    index := 0
    // Find good entry point so we don't accidently skip half the list.
    for i := t.maxLevel; i >= 0; i-- {
        if t.levels[i] != nil && t.levels[i].value.Compare(e) <= 0 {
            index = i
            break
        }
    }
    currentNode := t.levels[index]

    currCompare := currentNode.value.Compare(e)
    nextCompare := 0

    for {
        if currCompare == 0 {
            return currentNode, true
        }

        nextNode := currentNode.array[index].next
        if nextNode != nil {
            nextCompare = nextNode.value.Compare(e)
            currCompare = nextCompare
        }

        if nextNode != nil && nextCompare <= 0 {
            currentNode = nextNode
        } else {
            if index > 0 {
                index--
            } else {
                return nil, false
            }
        }
    }

    return nil, false

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

            if t.levels[i] == elem {
                t.levels[i] = next
                if next == nil {
                    // reduce the maximum entry position!
                    t.maxLevel = i-1
                }
            }
        }
    }
}

func (t *SkipList) PrettyPrint() {

    fmt.Printf("--> ")
    for i,l := range t.levels {
        next := "---"
        if l != nil {
            next = l.value.String()
        }
        fmt.Printf("[---|%v]", next)
        if i < len(t.levels)-1 {
            fmt.Printf(" --> ")
        }
    }
    fmt.Println("")

    node := t.levels[0]
    for node != nil {
        fmt.Printf("%v: ", node.value)
        for i,l := range node.array {

            prev := "---"
            if l.prev != nil {
                prev = l.prev.value.String()
            }
            next := "---"
            if l.next != nil {
                next = l.next.value.String()
            }

            fmt.Printf("[%v|%v]", prev, next)
            if i < len(node.array)-1 {
                fmt.Printf(" --> ")
            }

        }
        fmt.Printf("\n")
        node = node.array[0].next
    }
    fmt.Printf("\n")
}
