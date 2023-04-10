package main

import "github.com/gin-gonic/gin"

// 单个聊天室的JSON信息, 被ListRequestResp持有
type roomInfoJSON struct {
	Index int    `json:"id"`
	Title string `json:"title"`
}

// 用户访问http接口, 返回的数据
type ListRequestResp struct {
	Code      int             `json:"code"`
	Msg       string          `json:"msg"`
	RoomInfos []*roomInfoJSON `json:"room_infos"`
}

// 应该被go启动, 开启查询房间号的http服务, 有两个状态码 10000代表成功, 10001代表鉴权失败
func GoHttpService() {
	engine := gin.Default()
	engine.GET("/list", func(ctx *gin.Context) {
		username := ctx.Request.Header.Get("username")
		password := ctx.Request.Header.Get("password")
		println(username, password)
		if u, exists := Users[username]; exists && u.Password == password {
			resp := &ListRequestResp{}
			resp.Code = 10000
			resp.Msg = "success"
			for _, v := range Rooms {
				resp.RoomInfos = append(resp.RoomInfos, &roomInfoJSON{v.ID, v.Title})
			}
			ctx.JSON(200, resp)
		} else {
			ctx.JSON(200, &ListRequestResp{Code: 10001, Msg: "授权失败"})
			return
		}
	})
	engine.Run("localhost:8000")
}
