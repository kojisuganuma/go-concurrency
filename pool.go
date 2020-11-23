package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

const (
	numberOfExecute int = 10000
)

func main() {
	debug.SetGCPercent(-1) // GCが悪さをしてカウントが若干増えることがあるので無効化する
	var numberOfCreated int
	var m sync.Mutex

	// プール作成
	pool := &sync.Pool{
		New: func() interface{} {
			m.Lock()
			numberOfCreated++ // 実験のためリソースの新規作成回数をカウント
			m.Unlock()
			return &struct{}{}
		},
	}

	// 初期リソースの作成
	for i := 0; i < 1500; i++ {
		pool.Put(pool.New())
	}

	fmt.Println("")
	fmt.Println("===")
	fmt.Println("warm completed")

	var wg sync.WaitGroup
	wg.Add(numberOfExecute)
	for i := numberOfExecute; i > 0; i-- { // goroutineを10,000個生成
		go func() {
			defer wg.Done()
			s, ok := pool.Get().(*struct{}) // リソースの取得（枯渇時は内部でリソースが新規作成される）
			if !ok {
				panic("なんか変")
			}
			time.Sleep(1 * time.Millisecond) // 何らかの処理を模すため1ms待機
			pool.Put(s)                      // リソースの返却
		}()
	}
	wg.Wait()

	fmt.Printf("Number of created: %d\n", numberOfCreated)
}
