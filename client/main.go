package main

import "fmt"

var userId int     //用户的id
var userPwd string // 用户的密码
func main() {

	//接受用户的选择
	var key int
	//判断是否还继续显示菜单
	var loop = true

	for loop {
		fmt.Println("-------------------欢迎登陆多人聊天系统----------------------------------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择（1-3）")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			loop = false
		case 2:
			fmt.Println("注册用户")
			loop = false
		case 3:
			fmt.Println("退出系统")
			loop = false
		default:
			fmt.Println("您的输入有误，请重新输入")
		}
	}

	//根据用户输入，显示新的提示信息
	if key == 1 {
		//说明用户要登陆
		fmt.Println("请输入用户的id")
		fmt.Scanf("%d\n", &userId)
		fmt.Println("请输入用户的密码")
		fmt.Scanf("%s\n", &userPwd)
		//先把登录得函数写到另外一个文件,比如login.go
		login(userId, userPwd)
		//if err != nil {
		//	fmt.Println("登陆失败")
		//} else {
		//	fmt.Println("登陆成功")
		//}

	} else if key == 2 {
		fmt.Println("进行用户注册")
	}

}
