package skiplist

import (
    "fmt"
    "time"
    "testing"
    "math/rand"
    //"github.com/pkg/profile"
)


var g_maxN int = 1000000

//type Element struct {
//    E int
//}

type Element int

func (e Element) ExtractValue() float64 {
    return float64(e)
}
func (e Element) String() string {
    //return strconv.Itoa(e)
    return fmt.Sprintf("%03d", e)
}


// timeTrack will print out the number of nanoseconds since the start time divided by n
// Useful for printing out how long each iteration took in a benchmark
func timeTrack(start time.Time, n int, name string) {
    loopNS := time.Since(start).Nanoseconds() / int64(n)
    fmt.Printf("%s: %d\n", name, loopNS)
}


func TestBenchmarkInsert(t *testing.T) {
    list := New(0.01)

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element(g_maxN-i))
    }

    for i := 0; i < g_maxN; i++ {
        if _,ok := list.Find(Element(g_maxN-i)); !ok {
            t.Fail()
        }
    }

}
func TestBenchmarkWorstInsert(t *testing.T) {
    list := New(0.01)

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element(i))
    }

    for i := 0; i < g_maxN; i++ {
        if _,ok := list.Find(Element(i)); !ok {
            t.Fail()
        }
    }

}

func TestBenchmarkDelete(t *testing.T) {
    list := New(0.01)

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element(i))
    }

    for i := 0; i < g_maxN; i++ {
        list.Delete(Element(i))
    }

    if !list.isEmpty() {
        t.Fail()
    }
}
func TestBenchmarkWorstDelete(t *testing.T) {
    list := New(0.01)

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element(i))
    }

    for i := 0; i < g_maxN; i++ {
        list.Delete(Element(g_maxN-i-1))
    }
    if !list.isEmpty() {
        t.Fail()
    }
}

func TestFindGreaterOrEqual(t *testing.T) {
    list := New(0.01)

    for i := 0; i < g_maxN; i++ {
        if  i != 45 &&
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

    if e,ok := list.FindGreaterOrEqual(Element(44)); ok {
        if e.value.(Element) != 44 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element(45)); ok {
        if e.value.(Element) != 49 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element(47)); ok {
        if e.value.(Element) != 49 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element(6006)); ok {
        if e.value.(Element) != 6008 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element(6001)); ok {
        if e.value.(Element) != 6002 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element(6002)); ok {
        if e.value.(Element) != 6002 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

}

func TestInsertRandom(t *testing.T) {

    //defer profile.Start(profile.CPUProfile).Stop()
    list := New(0.01)

    rList := rand.Perm(g_maxN)

    for _,e := range rList {
        list.Insert(Element(e))
    }

    for _,e := range rList {
        if _,ok := list.Find(Element(e)); !ok {
            t.Fail()
        }
    }

    for _,e := range rList {
        list.Delete(Element(e))
    }

    if !list.isEmpty() {
        t.Fail()
    }
}
func TestDeleteRandom(t *testing.T) {

    //defer profile.Start(profile.CPUProfile).Stop()
    list := New(0.01)

    rList := rand.Perm(g_maxN)

    for _,e := range rList {
        list.Insert(Element(e))
    }

    for _,e := range rList {
        list.Delete(Element(e))
    }
}


func TestInsertDelete(t *testing.T) {

    //defer profile.Start(profile.CPUProfile).Stop()
    list := New(0.01)

    rList := rand.Perm(g_maxN)

    for _,e := range rList {
        list.Insert(Element(e))
    }

    for _,e := range rList {
        list.Delete(Element(e))
        list.Insert(Element(e))
    }
}

// Delete and Insert based on search:
// mtInserts: 476
// mtWorstInserts: 806
// mtAvgSearch: 510
// mtSearchEnd: 318
// mtDelete: 279
// mtWorstDelete: 421
// --- PASS: TestBenchmarkWorstDelete (1.15s)
// --- PASS: TestFind (1.33s)
// --- PASS: TestFindGreaterOrEqual (0.73s)
// --- PASS: TestDelete (1.10s)
// --- PASS: TestInsertRandom (3.91s)
// PASS
// ok   skiplist    12.916s


// Search finger introduced based on last insert/delete:
// mtInserts: 497
// mtWorstInserts: 823
// mtAvgSearch: 753
// mtSearchEnd: 82
// mtDelete: 512
// mtWorstDelete: 633
// --- PASS: TestBenchmarkWorstDelete (1.38s)
// --- PASS: TestFind (1.54s)
// --- PASS: TestFindGreaterOrEqual (0.81s)
// --- PASS: TestDelete (1.30s)
// --- PASS: TestInsertRandom (4.41s)
// ok   skiplist    14.486s

// No search finger but some other small improvements:
// mtInserts: 489
// mtWorstInserts: 450
// mtAvgSearch: 288
// mtSearchEnd: 188
// mtDelete: 126
// mtWorstDelete: 305
// --- PASS: TestBenchmarkWorstDelete (0.59s)
// --- PASS: TestFind (0.62s)
// --- PASS: TestFindGreaterOrEqual (0.35s)
// --- PASS: TestDelete (0.49s)
// --- PASS: TestInsertRandom (1.82s)
// --- PASS: TestDeleteRandom (3.43s)
// --- PASS: TestInsertDelete (4.93s)
// ok   skiplist    14.756s
