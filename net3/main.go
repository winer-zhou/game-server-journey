package main

import (
	"encoding/binary" // ← 新增：二进制编解码
	"fmt"
	"io" // ← 新增：io.ReadFull
	"net"
	"sync"
)

var (
	mu      sync.Mutex
	clients []net.Conn
)

// 发一条带长度前缀的消息
func send(conn net.Conn, msg string) {
	data := []byte(msg)                                   // ① 转成字节
	header := make([]byte, 4)                             // ② 造 4 字节头
	binary.BigEndian.PutUint32(header, uint32(len(data))) // ③ 把长度写进头

	conn.Write(header) // ④ 先发长度
	conn.Write(data)   // ⑤ 再发数据
}

func broadcast(msg string) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {
		send(c, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("有人连进来了:", conn.RemoteAddr())

	mu.Lock()
	clients = append(clients, conn)
	mu.Unlock()

	header := make([]byte, 4)

	for {
		// ⑥ 读满 4 字节的头
		if _, err := io.ReadFull(conn, header); err != nil {
			break
		}

		// ⑦ 解析出长度
		length := binary.BigEndian.Uint32(header)

		// ⑧ 读满 length 字节的数据
		data := make([]byte, length)
		if _, err := io.ReadFull(conn, data); err != nil {
			break
		}

		fmt.Println("收到", length, "字节:", string(data))
		broadcast(conn.RemoteAddr().String() + ": " + string(data))
	}

	mu.Lock()
	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mu.Unlock()

	fmt.Println("有人下线:", conn.RemoteAddr())
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动失败:", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务器启动，监听 8080 端口...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}
