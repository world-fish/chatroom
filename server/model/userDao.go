package model

import (
	"encoding/json"
	"fmt"
	"redigo-master/redis"
)

var (
	MyUserDao *UserDao
)

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	//通过 id 去 redis 查询这个用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil { //表示再 user  哈希中，没有找到对应的 id
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	//这里要把 res 反序列化成 User 实例
	json.Unmarshal([]byte(res), &user)

	return
}

func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//先从 UserDao 中取一个连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	} //这时证明这个用户id获取到了

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}
func (this *UserDao) Register(user *User) (err error) {
	//先从 UserDao 中取一个连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	//这是，说明 id 再 redis 里还没有，则可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	//入库
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err =", err)
		return
	}

	return

}
