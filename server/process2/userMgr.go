package process2

import "fmt"

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 因为 UserMgr 实例在服务器端只有一个，且在很多地方都会用到
// 将其定义为全局变量
var userMgr *UserMgr

// 完成对 userMgr 的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对 onlineUsers 的添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up

}

// 删除
func (this *UserMgr) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId)
}

// 返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据 id 返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up map[int]*UserProcess, err error) {
	//如何从map中取出一个值，带检测的方式
	_, ok := this.onlineUsers[userId]
	if !ok { //说明你要查找的这个用户，当前不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
