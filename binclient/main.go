package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// sendPacket：发送长度前缀包
func sendPacket(conn net.Conn, msg string) error {
	data := []byte(msg)
	var head [4]byte
	binary.BigEndian.PutUint32(head[:], uint32(len(data)))
	if _, err := conn.Write(head[:]); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

// readPacket：读一个完整包
func readPacket(r io.Reader) (string, error) {
	var head [4]byte
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return "", err
	}
	length := binary.BigEndian.Uint32(head[:])
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	fmt.Println("已连接（二进制协议）。输入命令：")

	// 后台收消息
	go func() {
		for {
			msg, err := readPacket(conn)
			if err != nil {
				fmt.Println("[连接断开]", err)
				os.Exit(0)
			}
			fmt.Print("[收到] " + msg)
		}
	}()

	// 主线程读键盘发送
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := sendPacket(conn, scanner.Text()); err != nil {
			fmt.Println("发送失败:", err)
			return
		}
	}
}
