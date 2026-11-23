package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	const N = 10000 // 并发连接数

	var wg sync.WaitGroup //等待所有gooutine完成
	successCount := 0
	var mu sync.Mutex //保护successCount

	start := time.Now()         // 记录开始时间
	var lastConnected time.Time // 最后一个连接建立成功的时刻

	for i := 0; i < N; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "localhost:8080")

			if err != nil {
				return
			}

			defer conn.Close()

			mu.Lock()
			successCount++
			lastConnected = time.Now()
			mu.Unlock()

			conn.Write([]byte(fmt.Sprintf("玩家%d 你好\n", id)))

			time.Sleep(10 * time.Second)
		}(i)
	}

	wg.Wait() //等待所有go程完成
	fmt.Printf("成功 %d / %d，全部建立耗时 %v\n",
		successCount, N, lastConnected.Sub(start))
	//	↑ 改成 successCount
	//
	// Sub = subtract 减法，算出两个时间点隔了多久
}
