package message

const (
	LoginMesType   = "LoginMes"
	LoginReMesType = "LoginReMes"
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
	Code  int    `json:"code"`  //返回状态码  500表示用户未注册  200表示登陆成功
	error string `json:"error"` //返回错误信息
}
