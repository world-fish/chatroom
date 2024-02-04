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
		//err = errors.New("read pkg header error")
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

func writePkg(conn net.Conn, data []byte) (err error) {
	var pkgLen uint32
	pkgLen = uint32(len(data)) //网络通信中经常要传递长度，而长度是非负的，因此选择无符号整数类型
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //网络通信中通常使用大端序BigEndian(高位在前)
	n, err := conn.Write(buf[0:4])               //发送数据的长度
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	//发送data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}

// 编写一个函数serverProcessLogin函数，专门处理登录请求
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	//核心代码
	//1.先从mes中取出mes.Data,并直接反序列成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return err
	}

	//2.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginMesType

	//再声明一个LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	//如果用户id=100 密码=123456，认为合法，否则不合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = 200
	} else {
		//不合法
		loginResMes.Code = 500 //500 状态码，表示用户不存在
		loginResMes.Error = "该用户不存在，请注册再使用"
	}

	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//4.将data赋值给resMes
	resMes.Data = string(data)

	//5.resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}

	//6.发送data,我们将其封装到writePkg函数中
	err = writePkg(conn, data)
	return

}

// 编写一个ServerProcessMes函数
// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登录逻辑
		serverProcessLogin(conn, mes)
	case message.RegisterMesType:
		//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
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
		//fmt.Println("mes=", mes)

		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}

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
