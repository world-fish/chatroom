package main

import (
	"chatroom/server/model"
	"fmt"
	"net"
	"time"
)

func process(conn net.Conn) {

	//这里需要延时关闭conn
	defer conn.Close()

	//这里调用总控，创建一个实例
	processor := &Processor{
		Conn: conn,
	}

	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务端通讯协程错误", err)
		return
	}
}

// 这里编写一个函数，完成对 UserDao 的初始化任务
func initUserDao() {
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	//服务器一开始，就初始化 redis 的 pool
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()

	//提示信息
	fmt.Println("服务器[新的结构]在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
	}

	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err)
		}

		//一旦链接成功，则启动一个携程和客户端保持通讯...
		go process(conn)
	}

}
