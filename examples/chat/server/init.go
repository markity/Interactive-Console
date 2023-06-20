package main

import "github.com/Allenxuxu/gev"

func init() {
	// 两个写死的账户
	Users["markity"] = &UserData{
		Username: "markity",
		Password: "123456",
	}
	Users["mary"] = &UserData{
		Username: "mary",
		Password: "56789",
	}

	// 两个写死的房间
	Rooms[1] = &RoomData{ID: 1, Title: "进来聊天啊", Messages: make([]*singleMessage, 0, 0)}
	Rooms[2] = &RoomData{ID: 2, Title: "进来聊骚啊", Messages: make([]*singleMessage, 0, 0)}

	Conns = make(map[int]*gev.Connection)
}
