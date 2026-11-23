package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const timeout = 300 * time.Second
const sessionTTL = 120 * time.Second // 会话保留 120 秒

// ════════════════════════════════════════════════════════════
// Client：新增 Send 字段（发送队列）
// ════════════════════════════════════════════════════════════
// // Session：断线期间保留的玩家数据
type Session struct {
	Name     string
	Room     string
	X        int
	Y        int
	LastSeen time.Time
}

type Client struct {
	Conn  net.Conn
	Room  string
	Name  string
	X     int
	Y     int
	Token string      // ← 新增：这个客户端的会话 token
	Send  chan string // ← 新增：这个客户端的发送队列
}

var (
	mu       sync.Mutex
	clients  []*Client
	sessions = make(map[string]*Session) //新增：断线期间保留的玩家数据
	wg       sync.WaitGroup              //等所有连接退出
)

// ════════════════════════════════════════════════════════════
// Push：把消息放进队列（不阻塞，队列满了就踢人）
// ════════════════════════════════════════════════════════════
func generateToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *Client) Push(msg string) {
	select {
	case c.Send <- msg:
		// 塞进去了，正常
	default:
		// 队列满了 = 这个客户端处理不过来 = 太慢
		fmt.Println("发送队列满，踢掉:", c.Name)
		c.Conn.Close()
	}
}

// ════════════════════════════════════════════════════════════
// WriteLoop：专门负责给这个客户端发消息
// ════════════════════════════════════════════════════════════
func (c *Client) WriteLoop() {
	for msg := range c.Send {
		if _, err := c.Conn.Write([]byte(msg)); err != nil {
			return // 连接断了，退出
		}
	}
}

// ════════════════════════════════════════════════════════════
// broadcastToRoom：只塞队列，不做网络 IO
// ════════════════════════════════════════════════════════════
func broadcastToRoom(room string, msg string) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {
		if c.Room == room {
			c.Push(msg) // 内存操作，瞬间完成
		}
	}
}

// ════════════════════════════════════════════════════════════
// handleConn
// ════════════════════════════════════════════════════════════
func handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defer wg.Done() // ← 结束时通知"我退了"

	client := &Client{
		Conn: conn,
		Send: make(chan string, 100), // 队列容量 100 条
	}

	mu.Lock()
	clients = append(clients, client)
	mu.Unlock()

	go client.WriteLoop() // 启动发送 goroutine

	// 因为下面的 scanner.Scan() 是阻塞的，没法自己感知 ctx
	// 只能靠"关掉 conn"让它返回
	connCtx, connCancel := context.WithCancel(ctx)
	defer connCancel()

	go func() {
		<-connCtx.Done()
		conn.Close() // 逼 Scan 返回
	}()

	// ── 欢迎 ──
	client.Push("欢迎！请输入: join 房间名 昵称\n")

	scanner := bufio.NewScanner(conn)
	conn.SetReadDeadline(time.Now().Add(timeout))

	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(timeout))
		line := scanner.Text()

		// ── join ──
		if strings.HasPrefix(line, "join ") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				client.Push("用法: join 房间名 昵称\n")
				continue
			}

			token := generateToken()

			mu.Lock()
			client.Room = parts[1]
			client.Name = parts[2]
			client.Token = token // ← 新增：生成会话 token

			// 创建会话
			sessions[token] = &Session{
				Name:     parts[2],
				Room:     parts[1],
				X:        0,
				Y:        0,
				LastSeen: time.Now(),
			}

			// 把房间里已有玩家的位置发给他
			for _, c := range clients {
				if c.Room == client.Room && c != client {
					client.Push(fmt.Sprintf("MOVE %s %d %d\n", c.Name, c.X, c.Y))
				}
			}
			mu.Unlock()

			client.Push("你进入了房间 " + parts[1] + "，昵称 " + parts[2] + "\n")
			client.Push("TOKEN " + token + "\n") // ← 把 token 发给客户端
			continue
		}

		// ── resume 断线重连 ──
		if strings.HasPrefix(line, "resume ") {
			token := strings.TrimSpace(strings.TrimPrefix(line, "resume "))

			mu.Lock()
			s, ok := sessions[token]
			if !ok {
				mu.Unlock()
				client.Push("会话已过期，请重新 join\n")
				continue
			}

			// 恢复会话
			client.Name = s.Name
			client.Room = s.Room
			client.X = s.X
			client.Y = s.Y
			client.Token = token
			s.LastSeen = time.Now()

			// 同步房间里其他人的位置
			for _, c := range clients {
				if c != client && c.Room == client.Room {
					client.Push(fmt.Sprintf("MOVE %s %d %d\n", c.Name, c.X, c.Y))
				}
			}
			mu.Unlock()

			// 广播"我回来了"（让别人看到我的位置）
			broadcastToRoom(client.Room, fmt.Sprintf("MOVE %s %d %d\n", client.Name, client.X, client.Y))
			client.Push("恢复成功: " + client.Name + " 在 " + client.Room + "\n")
			continue
		}

		// ── ping 心跳 ──
		if line == "ping" {
			continue
		}

		// ── 必须先进房间 ──
		if client.Room == "" {
			client.Push("请先输入: join 房间名 昵称\n")
			continue
		}

		// ── move ──
		if strings.HasPrefix(line, "move ") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				client.Push("用法: move x y\n")
				continue
			}

			x, err := strconv.Atoi(parts[1])
			if err != nil {
				client.Push("x 必须是数字\n")
				continue
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				client.Push("y 必须是数字\n")
				continue
			}

			mu.Lock()
			client.X = x
			client.Y = y
			mu.Unlock()

			broadcastToRoom(client.Room, fmt.Sprintf("MOVE %s %d %d\n", client.Name, x, y))
			continue
		}

		// ── 普通聊天 ──
		broadcastToRoom(client.Room, client.Name+": "+line+"\n")
	}

	// ══════ 断开时：保留会话，只更新最后活跃时间 ══════
	mu.Lock()
	if client.Token != "" {
		if s, ok := sessions[client.Token]; ok {
			s.LastSeen = time.Now()
		}
	}

	// 从在线列表移除（连接确实断了）
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()

	// 超时 / 断开
	if err := scanner.Err(); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("客户端超时:", conn.RemoteAddr())
		} else {
			fmt.Println("客户端断开:", conn.RemoteAddr())
		}
	}
}

// ════════════════════════════════════════════════════════════
// main
// ════════════════════════════════════════════════════════════
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ① 监听 Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\n收到退出信号,正在关闭...")
		cancel() // 通知所有人
	}()

	// ══════ 新增：ctx 取消时关闭 listener，让 Accept 返回 ══════
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				for token, s := range sessions {
					if time.Since(s.LastSeen) > sessionTTL {
						delete(sessions, token)
						fmt.Println("会话已过期:", s.Name)
					}
				}
				mu.Unlock()
			}
		}
	}()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动失败:", err)
		return
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	fmt.Println("服务器启动(channel 版），监听 8080 端口...,按 Ctrl+C 退出")

	// ② 循环接收连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				// 是我们自己关的，正常
			default:
				fmt.Println("Accept 出错:", err)
			}
			break
		}

		wg.Add(1) // 每接一个连接 +1
		go handleConn(ctx, conn)
	}

	// ④ 等所有连接退出
	fmt.Println("等待所有连接退出...")
	wg.Wait()
	fmt.Println("全部连接已关闭，服务器退出")
}
