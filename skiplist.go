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
    MAX_LEVEL = 25
)

type ListElement interface {
    ExtractValue() float64
    String() string
}

type Backtrack struct {
    node    *SkipListElement
    level   int
}

type SkipListElement struct {
    next       [MAX_LEVEL]*SkipListElement
    level       int
    key         float64
    value       ListElement
    prev        *SkipListElement
}

type SkipList struct {
    startLevels         [MAX_LEVEL]*SkipListElement
    endLevels           [MAX_LEVEL]*SkipListElement
    maxNewLevel         int
    maxLevel            int
    elementCount        int
    elementSum          float64
    eps                 float64
}

// Package initialization
func init() {
    seed := time.Now().UTC().UnixNano()
    seed = 1530734648380737920
    fmt.Printf("seed: %v\n", seed)
    rand.Seed(seed)
}

func New(eps float64) SkipList {
    list := SkipList{
        startLevels:        [MAX_LEVEL]*SkipListElement{},
        endLevels:          [MAX_LEVEL]*SkipListElement{},
        maxNewLevel:        MAX_LEVEL,
        maxLevel:           0,
        elementCount:       0,
        elementSum:         0.0,
        eps:                eps,
    }

    return list
}

func (t *SkipList)generateLevel(maxLevel int) int {
    level := 0
    // First we apply some mask which makes sure that we don't get a level
    // above our desired level. Then we find the first set bit.
    var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) -1)
    zeroes := bits.TrailingZeros64(x)
    if zeroes <= maxLevel {
        level = zeroes
    } else {
        level = maxLevel-1
    }

    return level
}

func (t *SkipList) isEmpty() bool {
    return t.startLevels[0] == nil
}

// returns: found element, backtracking list: Includes the elements from the entry point down to the element (or possible insertion position)!, ok, if an element was found
//func (t *SkipList) findExtended(key float64, findGreaterOrEqual bool) (foundElem *SkipListElement, ok bool) {
//
//
//    foundElem = nil
//    ok = false
//
//    if t.isEmpty() {
//        return
//    }
//
//    index := 0
//    var currentNode *SkipListElement = nil
//
//    // Find good entry point so we don't accidently skip half the list.
//    for i := t.maxLevel; i >= 0; i-- {
//        if t.startLevels[i] != nil && t.startLevels[i].key <= key {
//            index = i
//            break
//        }
//    }
//
//    currentNode = t.startLevels[index]
//    nextNode := currentNode
//
//    for {
//        if math.Abs(currentNode.key - key) <= t.eps {
//            foundElem = currentNode
//            ok = true
//            return
//        }
//
//        nextNode = currentNode.next[index]
//
//        // Which direction are we continuing next time?
//        if nextNode != nil && nextNode.key <= key {
//            // Go right
//            currentNode = nextNode
//        } else {
//            if index > 0 {
//
//                // Early exit
//                if currentNode.next[0] != nil && math.Abs(currentNode.next[0].key - key) <= t.eps {
//                    foundElem = currentNode.next[0]
//                    ok = true
//                    return
//                }
//
//                // Go down
//                index--
//            } else {
//                // Element is not found and we reached the bottom.
//                if findGreaterOrEqual {
//                    foundElem = nextNode
//                    ok = nextNode != nil
//                }
//                return
//
//            }
//        }
//    }
//
//    return
//}

func (t *SkipList) findExtended(key float64, findGreaterOrEqual bool) (foundElem *SkipListElement, ok bool) {


    foundElem = nil
    ok = false

    if t.isEmpty() {
        return
    }

    index := 0
    var currentNode *SkipListElement = nil

    // Find good entry point so we don't accidently skip half the list.
    for i := t.maxLevel; i >= 0; i-- {
        if t.startLevels[i] != nil && t.startLevels[i].key <= key {
            index = i
            break
        }
    }

    currentNode = t.startLevels[index]
    nextNode := currentNode

    for {
        if math.Abs(currentNode.key - key) <= t.eps {
            foundElem = currentNode
            ok = true
            return
        }

        nextNode = currentNode.next[index]

        // Which direction are we continuing next time?
        if nextNode != nil && nextNode.key <= key {
            // Go right
            currentNode = nextNode
        } else {
            if index > 0 {

                // Early exit
                if currentNode.next[0] != nil && math.Abs(currentNode.next[0].key - key) <= t.eps {
                    foundElem = currentNode.next[0]
                    ok = true
                    return
                }
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
    return t.findExtended(e.ExtractValue(), false)
}

func (t *SkipList) FindGreaterOrEqual(e ListElement) (*SkipListElement, bool) {
    return t.findExtended(e.ExtractValue(), true)
}

func (t *SkipList) Delete(e ListElement) {

   if t.isEmpty() {
        return
    }

    key := e.ExtractValue()

    //isFirst := false
    //isLast := false
    //if !t.isEmpty() {
    //    isFirst = math.Abs(key - t.startLevels[0].key) <= t.eps
    //    isLast  = math.Abs(key - t.endLevels[0].key  ) <= t.eps
    //}

    //normallyRemoved := false

    //if !isFirst && !isLast {

        //normallyRemoved = true

        index := 0

        // Find good entry point so we don't accidently skip half the list.
        for i := t.maxLevel; i >= 0; i-- {
            if t.startLevels[i] != nil && t.startLevels[i].key <= key {
                index = i
                break
            }
        }

        var currentNode *SkipListElement = nil
        nextNode := currentNode

        for {

            if currentNode == nil {
                nextNode = t.startLevels[index]
            } else {
                nextNode = currentNode.next[index]
            }

            // Found and remove!
            if nextNode != nil && math.Abs(nextNode.key - key) <= t.eps {


                if currentNode != nil {
                    currentNode.next[index] = nextNode.next[index]
                }


                if index == 0 && nextNode.next[index] != nil {
                    nextNode.next[index].prev = currentNode
                }

                // Link from start needs readjustments.
                if t.startLevels[index] == nextNode {
                    t.startLevels[index] = nextNode.next[index]
                    // This was our currently highest node!
                    if t.startLevels[index] == nil {
                        t.maxLevel = index -1
                    }
                }

                // Link from end needs readjustments.
                if nextNode.next[index] == nil {
                    t.endLevels[index] = currentNode
                }
                nextNode.next[index] = nil
            }

            if nextNode != nil && nextNode.key < key {
                // Go right
                currentNode = nextNode
            } else {
                // Go down
                index--
                if index < 0 {
                    break
                }
            }
        }
    //}

    // Where we have a left-most position that needs to be referenced!
    //for  i := level; i >= 0; i-- {
    //
    //    didSomething := false
    //
    //    if isFirst || normallyRemoved  {
    //
    //
    //        //if elem.next[i].prev == nil {
    //        if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
    //            if i == 0 && t.startLevels[i] != nil {
    //                t.startLevels[i].prev = elem
    //            }
    //            elem.next[i] = t.startLevels[i]
    //            t.startLevels[i] = elem
    //        }
    //
    //        // link the endLevels to this element!
    //        if elem.next[i] == nil {
    //            t.endLevels[i] = elem
    //        }
    //
    //        didSomething = true
    //    }
    //
    //    if isLast {
    //        // Places the element after the very last element on this level!
    //        // This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
    //        if !isFirst {
    //            if t.endLevels[i] != nil {
    //                t.endLevels[i].next[i] = elem
    //            }
    //            if i == 0 {
    //                elem.prev = t.endLevels[i]
    //            }
    //            t.endLevels[i] = elem
    //        }
    //
    //        // Link the startLevels to this element!
    //        //if elem.next[i].prev == nil {
    //        if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
    //            t.startLevels[i] = elem
    //        }
    //
    //        didSomething = true
    //    }
    //
    //    if !didSomething {
    //        break
    //    }
    //}

}

//func (t *SkipList) Delete(e ListElement) {
//
//    if elem,ok := t.Find(e); ok {
//        for i := elem.level; i >= 0; i-- {
//            prev := elem.next[i].prev
//            next := elem.next[i].next
//
//            if prev != nil {
//                prev.next[i].next = next
//            }
//            if next != nil {
//                next.next[i].prev = prev
//            }
//
//            if t.startLevels[i] == elem {
//                t.startLevels[i] = next
//                if next == nil {
//                    // reduce the maximum entry position!
//                    t.maxLevel = i-1
//                }
//            }
//            if t.endLevels[i] == elem {
//                t.endLevels[i] = prev
//            }
//        }
//        t.elementCount--
//        t.elementSum -= elem.key
//    }
//}

func (t *SkipList) Insert(e ListElement) {

    level := t.generateLevel(t.maxNewLevel)

    // Only grow the height of the skiplist by one at a time!
    if level > t.maxLevel+1 {
        level = t.maxLevel+1
    }

    elem  := &SkipListElement {
                next: [MAX_LEVEL]*SkipListElement{},
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

        index := 0
        if level > t.maxLevel {
            t.maxLevel = level
        }
        // Find good entry point so we don't accidently skip half the list.
        for i := t.maxLevel; i >= 0; i-- {
            if t.startLevels[i] != nil && t.startLevels[i].key <= elem.key || i <= level {
                index = i
                break
            }
        }

        var currentNode *SkipListElement = nil
        nextNode := t.startLevels[index]

        for {

            if currentNode == nil {
                nextNode = t.startLevels[index]
            } else {
                nextNode = currentNode.next[index]
            }

            // Connect node to next
            if index <= level && (nextNode == nil || nextNode.key >= elem.key) {
                elem.next[index] = nextNode
                if currentNode != nil {
                    currentNode.next[index] = elem
                }
                //elem.next[index].prev = currentNode
                //if nextNode != nil {
                //    nextNode.next[index].prev = elem
                //}
                if index == 0 {
                    elem.prev = currentNode
                    nextNode.prev = elem
                }
            }

            if nextNode != nil && nextNode.key < elem.key {
                // Go right
                currentNode = nextNode
            } else {
                // Go down
                index--
                if index < 0 {
                    break
                }
            }
        }
    }

    if level > t.maxLevel {
        t.maxLevel = level
    }

    // Where we have a left-most position that needs to be referenced!
    for  i := level; i >= 0; i-- {

        didSomething := false

        if newFirst || normallyInserted  {


            //if elem.next[i].prev == nil {
            if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
                if i == 0 && t.startLevels[i] != nil {
                    t.startLevels[i].prev = elem
                }
                elem.next[i] = t.startLevels[i]
                t.startLevels[i] = elem
            }

            // link the endLevels to this element!
            if elem.next[i] == nil {
                t.endLevels[i] = elem
            }

            didSomething = true
        }

        if newLast {
            // Places the element after the very last element on this level!
            // This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
            if !newFirst {
                if t.endLevels[i] != nil {
                    t.endLevels[i].next[i] = elem
                }
                if i == 0 {
                    elem.prev = t.endLevels[i]
                }
                t.endLevels[i] = elem
            }

            // Link the startLevels to this element!
            //if elem.next[i].prev == nil {
            if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
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
            l := node.next[i]

            prev := "---"
            //if l.prev != nil {
            //    prev = l.prev.value.String()
            //}
            next := "---"
            if l != nil {
                next = l.value.String()
            }

            fmt.Printf("[%v|%v]", prev, next)
            if i < node.level {
                fmt.Printf(" --> ")
            }

        }
        fmt.Printf("\n")
        node = node.next[0]
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
