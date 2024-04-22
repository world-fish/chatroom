package main

import (
	"fmt"
	"redigo-master/redis"
)

func main() {
	//链接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err=", err)
		return
	}
	defer conn.Close()
	fmt.Println()
	fmt.Println()

	//写入
	_, err = conn.Do("set", "name", "tjj")
	conn.Do("expire", "name", 10)
	if err != nil {
		fmt.Println("set err =", err)
		return
	}

	//读取数据
	fmt.Println("读取")
	r, err := redis.String(conn.Do("get", "name"))
	if err != nil {
		fmt.Println("set err =", err)
	}
	fmt.Println("操作ok", r)
}
