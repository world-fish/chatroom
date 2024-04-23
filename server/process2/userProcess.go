package process2

import (
	"chatroom/common/message"
	"chatroom/server/model"
	"chatroom/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段，表示该Conn是那个用户
	UserId int
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	//1.先从mes中取出mes.Data,并直接反序列成LoginMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return err
	}

	//2.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//我们需要到 redis 数据库去完成注册
	//1.使用 model.MyUserDao 到 redis 去验证
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误......"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
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
	//因为使用分层模式(mvc)，我们先拆功能键一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return

}

// 编写一个函数serverProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
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

	//我们需要到 redis 数据库去完成验证
	//使用 model.MyUserDao 到 redis 去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"

		}
		//不合法
	} else {
		loginResMes.Code = 200
		//这里，因为用户已经登陆成功，我们就把该登陆成功的用户放入到 userMgr 中
		//将登陆成功的用户的 userId 赋值给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//将当前在线用户的 id 放入到 loginResMes.UsersId 中
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)

		}

		fmt.Println(user, "登陆成功")
	}

	//如果用户id=100 密码=123456，认为合法，否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	//合法
	//	loginResMes.Code = 200
	//} else {
	//	//不合法
	//	loginResMes.Code = 500 //500 状态码，表示用户不存在
	//	loginResMes.Error = "该用户不存在，请注册再使用"
	//}

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
	//因为使用分层模式(mvc)，我们先拆功能键一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return

}
