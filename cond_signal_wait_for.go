package main

import (
    "sync"
    "fmt"
    "time"
)

func main() {
    var wg sync.WaitGroup
    var sharedRsc = make(map[string]interface{})

    wg.Add(2)
    m := sync.Mutex{}
    c := sync.NewCond(&m)
    go func() {
        // this go routine wait for changes to the sharedRsc
        c.L.Lock()
        for len(sharedRsc) == 0 {
	    fmt.Println("wait 1")
            c.Wait()
	    fmt.Println("arrive signal 1")
        }
        fmt.Println(sharedRsc["rsc1"])
        c.L.Unlock()
        wg.Done()
    }()

    go func() {
        // this go routine wait for changes to the sharedRsc
        c.L.Lock()
        for len(sharedRsc) == 0 {
	    fmt.Println("wait 2")
            c.Wait()
	    fmt.Println("arrive signal 2")
        }
        fmt.Println(sharedRsc["rsc2"])
        c.L.Unlock()
        wg.Done()
    }()

    time.Sleep(1 * time.Second)
    // this one writes changes to sharedRsc
    c.L.Lock()
    sharedRsc["rsc1"] = "foo"
    sharedRsc["rsc2"] = "bar"
    c.L.Unlock()
    time.Sleep(3 * time.Second)
    c.Signal()
    time.Sleep(3 * time.Second)
    c.Signal()
    wg.Wait()
}


