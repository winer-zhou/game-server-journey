package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	mu      sync.Mutex //锁:保护下面的clients
	clients []net.Conn //存所有在线的人
)

// 广播:把msg发给所有在线的人
func broadcast(msg string) {
	mu.Lock()
	defer mu.Unlock()

	for _, conn := range clients { //“_”表示忽略该值,声明的变量必须用
		conn.Write([]byte(msg))
	}
}

func handleConn(conn net.Conn) { // ① 把处理逻辑抽成函数
	defer conn.Close() // ② 这里用 defer 就安全了

	//上线
	mu.Lock()
	clients = append(clients, conn)
	mu.Unlock()

	scanner := bufio.NewScanner(conn) // 读取客户端发来的数据直到遇到换行符\n为止

	for scanner.Scan() { //只要还能读到行就一直循环
		line := scanner.Text()                                     //读取到的内容
		broadcast(conn.RemoteAddr().String() + ": " + line + "\n") //广播给所有人
	}

	//下线
	mu.Lock()
	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
		mu.Unlock()
	}
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
			fmt.Println("接受失败:", err)
			continue
		}

		go handleConn(conn) // ③ 关键：前面加个 go
	}
}
