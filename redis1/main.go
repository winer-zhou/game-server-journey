package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background() // 你刚学的 context，这里又用上了

	// 连接 Redis（返回的是连接池，不是单条连接）
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// PING 测试
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	fmt.Println("连接成功:", pong)
}
