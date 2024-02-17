package main

import (
	"chatroom/common/message"
	"chatroom/server/process2"
	"chatroom/server/utils"
	"fmt"
	"io"
	"net"
)

// Processor 先创建一个Processor的结构体
type Processor struct {
	Conn net.Conn
}

// 编写一个ServerProcessMes函数
// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登录逻辑
		//创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return

}

func (this *Processor) process2() (err error) {

	for {
		//这里我们将读取数据包，直接发嗯装成一个函数readPkg()，返回Message，Err
		//创建一个Transfer 实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}

		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		//fmt.Println("mes=", mes)

		err = this.ServerProcessMes(&mes)
		if err != nil {
			return err
		}

	}
}
