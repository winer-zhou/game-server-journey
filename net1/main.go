package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动失败", err)
		return
	}
	defer listener.Close()

	fmt.Println("服务器启动,监听8080端口...")

	for {
		conn, err := listener.Accept() //等待，有人进来才往下走,conn是连接对象的意思
		if err != nil {
			fmt.Println("接受连接失败:", err)
			continue
		}

		fmt.Println("有人连进来了:", conn.RemoteAddr())

		conn.Write([]byte("hello,我是服务器\n")) //发数据给客户端

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取失败", err)
			conn.Close()
			continue
		}

		fmt.Println("收到", n, "字节", string(buf[:n]))

		conn.Write([]byte("收到你的消息了\n"))
		conn.Close() //关闭连接
	}
}
