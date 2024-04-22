package model

type User struct {
	//为了保证序列化成功 用户信息的json和key 和 结构体的字段对应的tag名字一致
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
