package main

import "fmt"

// 写一个函数，完成登录校验
func login(userId int, userPwd string) (err error) {
	fmt.Printf("userId = %d userPwd = %s\n", userId, userPwd)
	return nil

}
