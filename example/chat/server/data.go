package main

// 用户的账号密码信息
type UserData struct {
	Username string
	Password string
	Logined  bool
}

// 单个聊天室的数据
type RoomData struct {
	ID       int
	Title    string
	Messages []*singleMessage
}
