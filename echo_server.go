package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func echo_worker(conn net.Conn) {
	defer conn.Close()
	// 如果用的时ipv4的tcp监听的话，这里返回的类型就是*net.TCPConn，这个类型实现了net.Conn接口
	// *net.TCPConn同样实现了io.Reader接口
	fmt.Printf("type of conn is %T, from %s\n", conn, conn.RemoteAddr().String())
	// content就是读到的内容
	content := make([]byte, 512)
	for {
		n, err := conn.Read(content)
		if err != nil {
			if err == io.EOF {
				//
				break
			} else {
				fmt.Printf("err when read: %v\n", err)
			}
		}
		// 读入了n个字节的内容, 原封不动地往回写
		n, err = conn.Write(content[0:n])
		if err != nil {
			fmt.Printf("err when writing: %v\n", err)
		} else {
			fmt.Printf("wrote %d bytes\n", n)
		}
	}
	fmt.Printf("conn from %s closed.\n", conn.RemoteAddr().String())
}

// 样式怎样搭建一个echo server的服务
func echo_server() {
	// 直接创建一个Listener
	var lstn net.Listener
	var err error
	lstn, err = net.Listen("tcp4", "127.0.0.1:8080") // 这里我们指定了要用ipv4的方式来监听，所以返回的Listener接口的真实类型时*net.TCPListener
	fmt.Printf("type of lstn = %T\n", lstn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for {
		// 接收连接
		var conn net.Conn
		conn, err = lstn.Accept()
		if err != nil {
			fmt.Printf("err at accept: %v\n", err)
			continue
		}
		go echo_worker(conn)
	}
}
