package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080") // ① 主动连接服务器
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	msg := "hello 长度前缀" // ② 要发的消息

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(msg))) // ③ 造头

	conn.Write(header)      // ④ 发长度
	conn.Write([]byte(msg)) // ⑤ 发数据

	fmt.Println("发送了", len(msg), "字节")
}
