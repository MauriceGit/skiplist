package skiplist

import (
    "fmt"
    "time"
    "testing"
    "math/rand"
    "github.com/pkg/profile"
)


var g_maxN int = 1000000

type Element struct {
    E int
}

func (e Element) Compare(e2 ListElement) int {

    if e.E < e2.(Element).E {
        return -1
    }
    if e.E == e2.(Element).E {
        return 0
    }
    return 1
}
func (e Element) String() string {
    //return strconv.Itoa(e.E)
    return fmt.Sprintf("%03d", e.E)
}


// timeTrack will print out the number of nanoseconds since the start time divided by n
// Useful for printing out how long each iteration took in a benchmark
func timeTrack(start time.Time, n int, name string) {
    loopNS := time.Since(start).Nanoseconds() / int64(n)
    fmt.Printf("%s: %d\n", name, loopNS)
}


func TestBenchmarkInsert(t *testing.T) {
    list := New()

    defer timeTrack(time.Now(), g_maxN, "mtInserts")

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{g_maxN-i})
    }
}
func TestBenchmarkWorstInsert(t *testing.T) {
    list := New()

    defer timeTrack(time.Now(), g_maxN, "mtWorstInserts")

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }
}
func TestBenchmarkAvgSearch(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }

    defer timeTrack(time.Now(), g_maxN, "mtAvgSearch")

    for i := 0; i < g_maxN; i++ {
        list.Find(Element{i})
    }
}
func TestBenchmarkSearchEnd(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }

    defer timeTrack(time.Now(), g_maxN, "mtSearchEnd")

    for i := 0; i < g_maxN; i++ {
        list.Find(Element{g_maxN-1})
    }
}
func TestBenchmarkDelete(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }

    defer timeTrack(time.Now(), g_maxN, "mtDelete")

    for i := 0; i < g_maxN; i++ {
        list.Delete(Element{i})
    }
}
func TestBenchmarkWorstDelete(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }

    defer timeTrack(time.Now(), g_maxN, "mtWorstDelete")

    for i := 0; i < g_maxN; i++ {
        list.Delete(Element{g_maxN-i})
    }
}

func TestFind(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }
    for i := 0; i < g_maxN; i++ {

        if e,ok := list.Find(Element{i}); !ok || e == nil {
            t.Fail()
        }
    }
}

func TestFindGreaterOrEqual(t *testing.T) {
    list := New()

    for i := 0; i < g_maxN; i++ {
        if  i != 45 &&
            i != 46 &&
            i != 47 &&
            i != 48 &&
            i != 6006 &&
            i != 6007 &&
            i != 6001 &&
            i != 6003 {
            list.Insert(Element{i})
        }
    }

    if e,ok := list.FindGreaterOrEqual(Element{44}); ok {
        if e.value.(Element).E != 44 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element{45}); ok {
        if e.value.(Element).E != 49 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element{47}); ok {
        if e.value.(Element).E != 49 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element{6006}); ok {
        if e.value.(Element).E != 6008 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element{6001}); ok {
        if e.value.(Element).E != 6002 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

    if e,ok := list.FindGreaterOrEqual(Element{6002}); ok {
        if e.value.(Element).E != 6002 {
            t.Fail()
        }
    } else {
        t.Fail()
    }

}

func TestDelete(t *testing.T) {


    list := New()

    for i := 0; i < g_maxN; i++ {
        list.Insert(Element{i})
    }

    //list.PrettyPrint()

    for i := 0; i < g_maxN; i++ {
        list.Delete(Element{i})
    }

    if !list.isEmpty() {
        t.Fail()
    }
}

func TestInsertRandom(t *testing.T) {

    defer profile.Start(profile.CPUProfile).Stop()
    list := New()

    rList := rand.Perm(g_maxN/2)
    for _,e := range rList {
        list.Insert(Element{e})
    }

    for _,e := range rList {
        if e2,ok := list.Find(Element{e}); !ok || e2 == nil {
            t.Fail()
        }
    }

    for _,e := range rList {
        list.Delete(Element{e})
    }

    if !list.isEmpty() {
        t.Fail()
    }
}
