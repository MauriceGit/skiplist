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

type Backtrack struct {
    node    *SkipListElement
    level   int
}

type SkipListElement struct {
    array       []SkipListPointer
    level       int
    value       ListElement


    // Unrolled Linked-List:
    //values      [5]ListElement
    //valueCount  int
    // Later experiment with caching the minimum and maximum element for faster comparison
    // when searching or inserting.
    //minIndex    int
    //maxIndex    int

}

type SkipList struct {
    levels              [25]*SkipListElement
    backtrack           []Backtrack
    lastBacktrackCount  int
    maxNewLevel         int
    maxLevel            int
}

// Package initialization
func init() {
    seed := time.Now().UTC().UnixNano()
    //seed = 1530040158521478743
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
    return SkipList{
        levels:             [25]*SkipListElement{},
        backtrack:          make([]Backtrack, 25),
        lastBacktrackCount: 0,
        maxNewLevel:        25,
        maxLevel:           0,
    }
}

func (t *SkipList) isEmpty() bool {
    return t.levels[0] == nil
}

// Uses the already filled backtrack slice to determine a better initial starting position to go down later.
// backtrack slice HAS to be filled beforehand!
// Returns a Skip list element and the index where a normal search should start!
// One big advantage is, that we can skip all level marching and will just go to the top of a node directly!!!
func (t *SkipList) useSearchFingerForEntry(e ListElement) (*SkipListElement, int, bool) {

    // Update: Use the existing backtrack-Slice and go up backwards.
    // When we are right of the element:
    //      Maybe, for each element, check the .prev. If it is <= our element, go there instead and Stop!
    // When we are left of the element: Go up then right as much as possible. As soon as .next is >= our element, Stop!

    node  := t.backtrack[t.lastBacktrackCount-1].node
    level := 0

    goLeft := e.Compare(node.value) < 0

    for i := t.lastBacktrackCount-1; i >= 0; i-- {

        node  = t.backtrack[i].node
        level = t.backtrack[i].level

        // We found a good starting point for either direction :)
        if goLeft {
            // No way left or we passed our element.
            if node.array[level].prev != nil {
                if e.Compare(node.array[level].prev.value) >= 0 {
                    return node.array[level].prev, level, true
                }
            } else {
                return node, level, false
            }
        } else {
            // No way left or we passed our element.
            if node.array[level].next == nil || e.Compare(node.array[level].next.value) <= 0 {
                return node, level, true
            }
        }
    }

    // We are probably back at the top, meaning - The element we look for is most likely not in our tree...
    // No, definitely not.
    return node, level, false
}

// returns: found element, backtracking list: Includes the elements from the entry point down to the element (or possible insertion position)!, ok, if an element was found
func (t *SkipList) findExtended(e ListElement, findGreaterOrEqual bool, createBackTrack bool) (*SkipListElement, *[]Backtrack, int, bool) {
    if t.isEmpty() {
        return nil, nil, 0, false
    }

    btCount := 0

    index := 0
    var currentNode *SkipListElement = nil

    useSearchFinger := false

    // Use the search finger for a good starting point entry.
    // We can't use the backtrack and create one at the same time (especially because we go wrong paths).
    if t.lastBacktrackCount >= 1 && !createBackTrack {
        n, i, ok := t.useSearchFingerForEntry(e)
        // We know, that the element we look for is NOT in in the skiplist.
        // We can only return early here because we will not use this now for a delete or an insert. And find is finished.
        if ok {
            currentNode = n
            index = i
            useSearchFinger = true
        }
    }

    if !useSearchFinger {
        // Find good entry point so we don't accidently skip half the list.
        for i := t.maxLevel; i >= 0; i-- {
            if t.levels[i] != nil && t.levels[i].value.Compare(e) <= 0 {
                index = i
                break
            }
        }
        currentNode = t.levels[index]
    }

    currCompare := currentNode.value.Compare(e)
    nextCompare := 0

    for {
        if currCompare == 0 {
            return currentNode, &t.backtrack, btCount, true
        }

        nextNode := currentNode.array[index].next
        if nextNode != nil {
            nextCompare = nextNode.value.Compare(e)
            currCompare = nextCompare
        }

        // Which direction are we continuing next time?
        if nextNode != nil && nextCompare <= 0 {
            // Go right
            currentNode = nextNode
        } else {
            if createBackTrack {
                t.backtrack[btCount].node = currentNode
                t.backtrack[btCount].level = index
                btCount++
            }
            if index > 0 {
                // Go down
                index--
            } else {
                // Element is not found and we reached the bottom.
                if findGreaterOrEqual {
                    return nextNode, &t.backtrack, btCount, nextNode != nil
                } else {
                    return nil, &t.backtrack, btCount, false
                }
            }
        }
    }

    return nil, nil, 0, false
}

func (t *SkipList) Find(e ListElement) (*SkipListElement, bool) {
    l, _, _, ok := t.findExtended(e, false, false)
    return l, ok
}

func (t *SkipList) FindGreaterOrEqual(e ListElement) (*SkipListElement, bool) {
    l, _, _, ok := t.findExtended(e, true, false)
    return l, ok
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

func (t *SkipList) Insert(e ListElement) {

    level := generateLevel(t.maxNewLevel)
    elem  := &SkipListElement{
                array: make([]SkipListPointer, level+1, level+1),
                level: level,
                value: e,
            }

    newFirst := true
    if !t.isEmpty() {
        newFirst = t.levels[0].value.Compare(e) > 0
    }

    // Insertion using Find()
    if !newFirst {

        // Using the search-finger approach, if


        // Search for e down to level 1. It will not find anything, but will return a backtrack for insertion.
        _, backtrack, btCount, _ := t.findExtended(e, true, true)
        // So we can use this backtrack the next time we look for an insertion position!
        t.lastBacktrackCount = btCount

        for i := btCount-1; i >= 0; i-- {

            bt := (*backtrack)[i]

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
