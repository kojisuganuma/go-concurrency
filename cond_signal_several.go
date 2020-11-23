package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    wg := new(sync.WaitGroup)

    l := new(sync.Mutex)
    c := sync.NewCond(l)
    for i := 0; i < 10; i++ {
	wg.Add(1)
        go func(i int) {
            fmt.Printf("waiting %d\n", i)
            c.L.Lock()
            defer c.L.Unlock()
	    defer wg.Done()
            c.Wait()
            fmt.Printf("go %d\n", i)
        }(i)
    }

    for i := 0; i < 5; i++ {
        time.Sleep(1 * time.Second)
        c.Signal()
    }
    time.Sleep(1 * time.Second)
    for i := 0; i < 5; i++ {
        time.Sleep(2 * time.Second)
        c.Signal()
    }

    wg.Wait()
}

