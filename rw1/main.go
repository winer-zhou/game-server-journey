package main

import (
	"fmt"
	"sync" // synchronize：同步工具（WaitGroup、Mutex、RWMutex）
	"time" // time：时间（Sleep、Now、Since）
)

func main() {
	// ══════════════════════════════════════════════════════════
	// 实验 1：普通 Mutex —— 两个"读"也要排队
	// 预期总耗时：4 秒（串行）
	// ══════════════════════════════════════════════════════════
	fmt.Println("=== 实验 1:Mutex 两个读 ===")

	t1 := time.Now() // Now()：记录开始时间

	var mu1 sync.Mutex     // Mutex：普通互斥锁（一次只准一个人进）
	var wg1 sync.WaitGroup // WaitGroup：等两个 goroutine 都做完

	for i := 1; i <= 2; i++ { // 循环 2 次，开 2 个 goroutine
		wg1.Add(1) // 计数 +1

		go func(n int) { // go：开 goroutine；n 是传进去的编号（1 或 2）
			defer wg1.Done() // defer：函数结束时计数 -1

			mu1.Lock()                   // 上锁：独占
			fmt.Printf("  读 %d 开始\n", n) // 我开始读了
			time.Sleep(2 * time.Second)  // 模拟"读了 2 秒"（占着锁）
			fmt.Printf("  读 %d 结束\n", n) // 我读完了
			mu1.Unlock()                 // 解锁

		}(i) // 把 i 作为 n 传进去
	}

	wg1.Wait() // 等两个都完成

	fmt.Printf("→ 总耗时: %v\n\n", time.Since(t1).Round(time.Second))
	// Since(t1) = 从 t1 到现在过了多久
	// Round(time.Second) = 四舍五入到秒（好看）

	// ══════════════════════════════════════════════════════════
	// 实验 2：RWMutex —— 两个"读"可以并行
	// 预期总耗时：2 秒（并行！）
	// ══════════════════════════════════════════════════════════
	fmt.Println("=== 实验 2:RWMutex 两个读 ===")

	t2 := time.Now()

	var mu2 sync.RWMutex // RWMutex：读写锁（读共享，写独占）
	var wg2 sync.WaitGroup

	for i := 1; i <= 2; i++ {
		wg2.Add(1)

		go func(n int) {
			defer wg2.Done()

			mu2.RLock() // RLock = Read Lock，读锁
			// 关键：多个 goroutine 可以同时拿到读锁！
			fmt.Printf("  读 %d 开始\n", n)
			time.Sleep(2 * time.Second)
			fmt.Printf("  读 %d 结束\n", n)
			mu2.RUnlock() // RUnlock = 释放读锁

		}(i)
	}

	wg2.Wait()

	fmt.Printf("→ 总耗时: %v\n\n", time.Since(t2).Round(time.Second))

	// ══════════════════════════════════════════════════════════
	// 实验 3：RWMutex —— 一个读 + 一个写（互斥）
	// 预期总耗时：4 秒（读写不能同时）
	// ══════════════════════════════════════════════════════════
	fmt.Println("=== 实验 3:RWMutex 一个读一个写 ===")

	t3 := time.Now()

	var mu3 sync.RWMutex
	var wg3 sync.WaitGroup

	// ── 第一个 goroutine：读 ──
	wg3.Add(1)
	go func() {
		defer wg3.Done()

		mu3.RLock() // 拿读锁
		fmt.Println("  读 开始")
		time.Sleep(2 * time.Second)
		fmt.Println("  读 结束")
		mu3.RUnlock() // 放开读锁
	}()

	// ── 第二个 goroutine：写 ──
	wg3.Add(1)
	go func() {
		defer wg3.Done()

		mu3.Lock() // 拿写锁（要等读锁全放开）
		fmt.Println("  写 开始")
		time.Sleep(2 * time.Second)
		fmt.Println("  写 结束")
		mu3.Unlock() // 放开写锁
	}()

	wg3.Wait()

	fmt.Printf("→ 总耗时: %v\n", time.Since(t3).Round(time.Second))
}
