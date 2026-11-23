package main

// ============================================================
// ① import —— 要用哪些工具箱
// ============================================================
// bufio   = buffered I/O，带缓冲的读写（按行读用它）
// fmt     = format，格式化打印
// net     = network，网络（监听、连接、收发）
// strings = 字符串工具（判断前缀、去掉前缀）
// sync    = synchronize，同步（互斥锁）
import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================
// ② Client —— 一个玩家长什么样
// ============================================================
// 把一个玩家的"连接"和"房间名"打包在一起。
// 以后要加昵称、血量、坐标，都往这个结构体里加字段即可。
type Client struct {
	Conn net.Conn // 怎么联系他
	Room string   // 他在哪个房间（空字符串 = 还没进房间）
	Name string   //昵称
	X    int
	Y    int
}

const timeout = 10 * time.Second // 连接超时时间

// ============================================================
// ③ 全局变量 —— 整个服务器的"记忆"
// ============================================================
// clients —— 所有在线玩家的列表（相当于店里的客户名册）
// mu      —— 保护 clients 的锁（同一时间只能一个人改）
//
// 为什么用 []*Client（指针切片）？
// 因为我们要修改 client.Room，值传递改不到原件，必须用指针。
var (
	mu      sync.Mutex
	clients []*Client
)

// ============================================================
// ④ broadcastToRoom —— 群发功能
// ============================================================
// 遍历所有在线玩家，只发给"房间名相同"的人。
// 这就是"房间隔离"的实现。
//
// 以后想改成"只发给附近的人"（游戏里的 AOI），
// 只需要改 if 那一行的判断条件。
func broadcastToRoom(room string, msg string) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {
		if c.Room == room {
			c.Conn.Write([]byte(msg))
		}
	}
}

// ============================================================
// ⑤ handleConn —— 一个玩家从上线到下线的一生
// ============================================================
// 分三个阶段：
//
//	上线 —— 把玩家加进 clients
//	循环 —— 不断读消息，判断是"命令"还是"普通聊天"
//	下线 —— 从 clients 里删掉自己
//
// 以后要加新命令（move / attack / buy），
// 就在循环里加一个 if 分支。
func handleConn(conn net.Conn) {
	defer conn.Close()

	// ── 阶段 1：上线 ──
	client := &Client{Conn: conn, Room: "", Name: "", X: 0, Y: 0}

	mu.Lock()
	clients = append(clients, client)
	mu.Unlock()

	// ── 阶段 2：循环读消息 ──
	conn.Write([]byte("欢迎! 请输入:join房间名 昵称 例如:join game1 zhou\n"))
	//怎么用建房间
	conn.Write([]byte("move 100 200\n"))

	scanner := bufio.NewScanner(conn)

	conn.SetReadDeadline(time.Now().Add(timeout)) // 设置读取超时时间

	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(timeout)) // 每次成功读取后重置超时时间
		line := scanner.Text()

		// 命令：join 房间名
		if strings.HasPrefix(line, "join ") {
			parts := strings.Fields(line)
			// Fields：按空格切分成切片
			// "join game1 zhou" → ["join", "game1", "zhou"]

			if len(parts) < 3 {
				conn.Write([]byte("用法: join 房间名 昵称\n"))
				continue
			}

			mu.Lock()
			client.Room = parts[1]
			client.Name = parts[2]
			// 把房间里已有玩家的位置，一个个发给新来的人
			for _, c := range clients {
				if c.Room == client.Room && c != client {
					conn.Write([]byte(fmt.Sprintf("MOVE %s %d %d\n", c.Name, c.X, c.Y)))
				}
			}

			mu.Unlock()

			conn.Write([]byte("你进入了房间 " + parts[1] + "，昵称 " + parts[2] + "\n"))
			continue
		}

		if line == "ping" {
			continue // 客户端发 ping 就不回 pong，客户端自己处理
		}

		// 还没进房间就想说话
		if client.Room == "" {
			conn.Write([]byte("请先输入 join 房间名\n"))
			continue
		}

		// ── 命令 2：move x y  ← 新增 ──
		if strings.HasPrefix(line, "move ") {
			parts := strings.Fields(line)
			// "move 100 200" → ["move", "100", "200"]

			if len(parts) < 3 {
				conn.Write([]byte("用法: move x y\n"))
				continue
			}

			// 把 x、y 字符串转成数字，就是把玩家输入的"文字"变成"真正的数字"
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				conn.Write([]byte("x 必须是数字\n"))
				continue
			}

			y, err := strconv.Atoi(parts[2])
			if err != nil {
				conn.Write([]byte("y 必须是数字\n"))
				continue
			}

			// 更新自己的位置
			mu.Lock()
			client.X = x
			client.Y = y
			mu.Unlock()

			// 广播给同房间所有人
			broadcastToRoom(client.Room, fmt.Sprintf("MOVE %s %d %d\n", client.Name, x, y))
			continue
		}
		// ── 普通聊天 ──
		broadcastToRoom(client.Room, client.Name+": "+line+"\n")
	}

	//看看是不是超时断的
	if err := scanner.Err(); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("客户端超时:", conn.RemoteAddr())
		} else {
			fmt.Println("客户端断开:", conn.RemoteAddr(), err)
		}
	}

	// ── 阶段 3：下线 ──
	mu.Lock()
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()
}

// ============================================================
// ⑥ main —— 接客
// ============================================================
// 只干两件事：
//  1. 监听端口（开门营业）
//  2. 无限循环 Accept（等客人），每来一个就开 goroutine 服务他
//
// 所以 main 永远不会被某一个玩家卡住。
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动失败:", err)
		return
	}
	defer listener.Close()

	fmt.Println("房间版服务器启动，监听 8080 端口...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}
