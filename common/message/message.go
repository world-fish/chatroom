package message

import (
	"chatroom/server/model"
)

const (
	LoginMesType       = "LoginMes"
	LoginResMesType    = "LoginReMes"
	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// 定义两个消息，后面需要再增加
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"`    //返回状态码  500表示用户未注册  200表示登陆成功
	UsersId []int  `json:"usersId"` //增加字段，保存用户id的切片
	Error   string `json:"error"`   //返回错误信息
}

type RegisterMes struct {
	User model.User `json:"user"`
}
type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码  400表示该用户已经占用  200表示登陆成功
	Error string `json:"error"` //返回错误信息
}
