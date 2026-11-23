package main

import (
	"context"
	"fmt"
	"time"
)

// 模拟一个耗时操作，使用 context 来处理超时
// 感知外部的取消信号
func slowOperation(ctx context.Context) (string, error) {
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done(): // 检查 context 是否已取消
			return "", ctx.Err() // 返回 context 的错误信息

		case <-time.After(1 * time.Second): // 模拟每次操作耗时 1 秒
			fmt.Println("第%d次操作", i)
		}
	}
	return "查询结果", nil
}

// 主函数，演示如何使用带有超时的 context 来调用耗时操作
func main() {
	// 创建一个带有超时的 context，超时时间为 3 秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 确保在函数退出时释放资源

	result, err := slowOperation(ctx) // 调用耗时操作
	if err != nil {
		fmt.Println("操作失败:", err) // 如果操作失败，打印错误信息
		return
	}
	fmt.Println("操作成功:", result) // 如果操作成功，打印结果
}
