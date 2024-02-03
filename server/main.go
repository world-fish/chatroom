package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {

	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	_, err = conn.Read(buf[:4])
	if err != nil {
		fmt.Println("conn.Read err=", err)
		//err = errors.New("read pkg hander error")
		return
	}

	//根据buf[:4] 转换成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//fmt.Println("conn.Read err=", err)
		err = errors.New("read pkg body error")
		return
	}

	//把pkgLen反序列化成->message.Message
	//技术是一层窗户纸
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}
func process(conn net.Conn) {

	//这里需要延时关闭conn
	defer conn.Close()

	for {
		mes, err := readPkg(conn)

		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出")
				return
			} else {
				fmt.Println("readPkg err=", err)
				return
			}
		}

		fmt.Println("mes=", mes)
	}

	//buf := make([]byte, 8096)

	//循环的读取客户端发送的信息
	//for {
	//
	//	//这里我们将读取的数据包，直接封装成一个函数readPkg()，返回Message，Err
	//
	//	fmt.Println("读取客户端发送的数据...")
	//	_, err := conn.Read(buf[:4])
	//	if err != nil {
	//		fmt.Println("conn.Read err=", err)
	//		return
	//	}
	//	fmt.Println("读到的buf=", buf[:4])
	//}

}

func main() {

	//提示信息
	fmt.Println("服务器在8889端口监听...")
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
