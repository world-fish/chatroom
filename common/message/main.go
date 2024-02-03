package message

const (
	LoginMesType   = "LoginMes"
	LoginReMesType = "LoginReMes"
)

type Message struct {
	Type string
	Data string
}

// 定义两个消息，后面需要再增加
type LoginMes struct {
	UserId   int    //用户id
	UserPwd  string //用户密码
	UserName string //用户名
}

type LoginResMes struct {
	Code  int    //返回状态码    500表示用户未注册  200表示登陆成功
	error string //返回错误信息
}
