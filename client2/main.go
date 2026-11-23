package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Pos = position 位置：一个坐标
type Pos struct {
	X int
	Y int
}

// Pos = position 位置：一个坐标
var (
	mu    sync.Mutex              // 互斥锁，用于保护共享资源
	peers = make(map[string]*Pos) // 存储在线玩家的坐标信息
)

func main() {
	// 监听 TCP 端口
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	fmt.Println("已连接。先输入:join 房间名 昵称")

	go receive(conn) // 接收消息的 goroutine

	go receive(conn)
	go heartbeat(conn) // 心跳 goroutine

	// 读取用户键盘输入的消息并发送给服务器
	scanner := bufio.NewScanner(os.Stdin) // 读取用户输入的消息
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n")) // 将消息发送给服务器
	}
}

// 接收服务器发送的消息
func receive(conn net.Conn) {
	scanner := bufio.NewScanner(conn) // 读取服务器发送的消息
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "MOVE") {
			// 解析坐标信息
			parts := strings.Fields(line)
			if len(parts) < 4 {
				continue
			}

			x, err1 := strconv.Atoi(parts[2])
			y, err2 := strconv.Atoi(parts[3])
			if err1 != nil || err2 != nil {
				continue
			}
			mu.Lock()
			peers[parts[1]] = &Pos{X: x, Y: y}
			mu.Unlock()

			drawMap() // 重画地图
			continue
		}
		fmt.Println(line) // 打印服务器发送的消息
	}
}

// 地图的宽高（字符格数）
const mapW = 40
const mapH = 8

// drawMap：清掉旧画面，重画一张地图
func drawMap() {
	mu.Lock()
	defer mu.Unlock()

	// ① 保存光标位置 —— 这样重画完能回到你输入的地方
	fmt.Print("\033[H\033[2J")
	// \033[s = 保存光标
	// \033[H = 光标移到左上角（1,1）

	// ② 画上边框
	fmt.Println("+" + strings.Repeat("-", mapW) + "+\033[K")
	// Repeat = 重复：把 "-" 重复 40 次

	// ③ 一行一行画
	for y := 0; y < mapH; y++ {
		row := make([]byte, mapW) // 一行 40 个格子
		for i := range row {
			row[i] = ' ' // 先全部填空格
		}

		// 把这一行上的玩家放上去
		for _, p := range peers {
			// 世界坐标(0~1000) → 屏幕格数(0~39)
			cx := p.X * (mapW - 1) / 1000
			cy := p.Y * (mapH - 1) / 1000

			if cy == y && cx >= 0 && cx < mapW {
				row[cx] = '@' // 画一个 @ 代表玩家
			}
		}

		fmt.Printf("|%s|\033[K\n", row) // %s 可以直接打印 []byte
	}

	// ④ 画下边框
	fmt.Println("+" + strings.Repeat("-", mapW) + "\033[K")

	// ⑤ 列出所有人和坐标
	info := ""
	for name, p := range peers {
		info += fmt.Sprintf("%s(%d,%d) ", name, p.X, p.Y)
	}
	fmt.Println(info + "\033[K")

	// ⑥ 恢复光标 —— 回到你原来输入的位置
	fmt.Printf("\033[%d;1H", mapH+6)
}

// heartbeat：每 5 秒发一次 ping，告诉服务器"我还活着"
func heartbeat(conn net.Conn) {
	ticker := time.NewTicker(5 * time.Second) // 定时器：每 5 秒触发一次
	defer ticker.Stop()                       // 函数结束时停止

	for range ticker.C { // 每 5 秒转一圈
		conn.Write([]byte("ping\n")) // 发心跳包
	}
}
