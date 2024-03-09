package process

import (
	"chatroom/common/message"
	"chatroom/server/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//暂时不需要字段...
}

// 给关联一个用户登录的方法
// 写一个函数，完成登录校验
func (this *UserProcess) Login(userId int, userPwd string) (err error) {

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

	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.将loginMes序列化
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

	//fmt.Printf("客户端，发送消息的长度=%d 内容=%s\n", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	//休眠2s
	//time.Sleep(2 * time.Second)
	//fmt.Println("休眠2s...")

	//这里还需要处理服务器端返回的消息
	//创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg() //mes就是
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	//将mes的Data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//fmt.Println("登陆成功")

		//这里我们还需要再客户端启动一个协程
		//该协程保持和服务端的通讯
		//如果服务器有数据推送给客户端 则接收并显示在客户端的终端
		go serverProcessMes(conn)

		//1.显示我们的登陆成功的菜单[循环]..
		for {
			ShowMenu()
		}
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}

	return

}
