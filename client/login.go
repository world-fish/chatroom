package main

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 写一个函数，完成登录校验
func login(userId int, userPwd string) (err error) {

	//下一个就要开始定协议...
	//fmt.Printf("userId = %d userPwd = %s\n", userId, userPwd)
	//return nil

	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送消息给服务
	var mes message.Message
	mes.Type = message.LoginMesType

	//3.创建一个LoginMes结构其
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.讲loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//5.把data赋值给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7.到这个时候 data就是我们要发送的信息
	//7.1 先把data的长度发送给服务器
	//先获取到data的长度->转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data)) //网络通信中经常要传递长度，而长度是非负的，因此选择无符号整数类型
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //网络通信中通常使用大端序BigEndian(高位在前)
	n, err := conn.Write(buf[0:4])               //发送数据的长度
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	fmt.Printf("客户端，发送消息的长度=%d 内容=%s\n", len(data), string(data))
	return

}
