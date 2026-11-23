package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background() // 创建一个 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379"}) // Redis 服务器地址
	defer rdb.Close() // 确保在程序结束时关闭客户端

	const key = "game:leaderboard" // 定义排行榜的键名

	// 清空排行榜
	rdb.Del(ctx, key)

	// 定义玩家及其分数
	players := map[string]float64{
		"zhou": 1500,
		"小明":   2300,
		"大佬":   3800,
		"新手":   200,
		"路人":   900,
	}

	// 将玩家及其分数添加到排行榜
	for name, score := range players {
		rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: name})
	}
	fmt.Println("已写入5个玩家")

	// ③ 给 zhou 加分 500（击杀奖励）
	rdb.ZIncrBy(ctx, key, 500, "zhou")
	fmt.Println("zhou 击杀 +500 分")

	// ④ 查前 3 名
	fmt.Println("\n=== 排行榜前 3 ===")
	result, err := rdb.ZRevRangeWithScores(ctx, key, 0, 2).Result()
	if err != nil {
		fmt.Println("查询失败:", err)
		return
	}
	for i, z := range result {
		fmt.Printf("%d. %s — %.0f 分\n", i+1, z.Member, z.Score)
	}

	// ⑤ 查 zhou 的排名
	rank, err := rdb.ZRevRank(ctx, key, "zhou").Result()
	if err != nil {
		fmt.Println("查排名失败:", err)
	} else {
		fmt.Printf("\nzhou 当前排第 %d 名\n", rank+1) // ZRevRank 从 0 开始，所以要 +1
	}
}
