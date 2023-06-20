package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/gev/plugins/websocket/ws/util"
)

var (
	KeyRequestHeader = "request_header" // http.Header
	KeyUri           = "uri"            // uri
	KeyRoomIndex     = "room_idx"       // int
	KeyUsername      = "username"       // string
	KeyConnID        = "conn_id"        //int
	KeyAuthed        = "authed"
)

type ChatServer struct{}

func (s *ChatServer) OnConnect(c *gev.Connection) {
	log.Println("OnConnect: ", c.PeerAddr())
}

func (s *ChatServer) OnMessage(c *gev.Connection, data []byte) (messageType ws.MessageType, out []byte) {
	log.Printf("OnMessage: %s\n", data)

	roomIdx_, _ := c.Get(KeyRoomIndex)
	roomIdx := roomIdx_.(int)
	username_, _ := c.Get(KeyUsername)
	username := username_.(string)

	switch DetectMessageType(data) {
	case TypeClientFetchMessages:
		room := Rooms[roomIdx]
		var s ClientFetchMessages
		_ = json.Unmarshal(data, &s)
		if s.IndexBefore < 0 || s.IndexBefore > len(room.Messages) {
			c.Close()
			return
		}

		var start int
		var end int
		var length = len(room.Messages)
		if s.IndexBefore == 0 {
			// 暂定为10条吧
			start = length - 10
			if start < 0 {
				start = 0
			}
			end = start + 10
			if end > length {
				end = length
			}
		} else {
			start = s.IndexBefore - 10
			if start < 0 {
				start = 0
			}
			end = start + 10
			if end > len(room.Messages) {
				end = len(room.Messages)
			}

		}
		fmt.Printf("start = %v, end = %v, len = %v", start, end, len(room.Messages))
		m := make([]*singleMessage, 0, 10)
		for ; start != end; start++ {
			m = append(m, room.Messages[start])
		}

		messageType = ws.MessageBinary
		fmt.Println(m)
		fmt.Println(string(MustMarshalToBytes(NewServerRespMessages(s.IndexBefore, m))))
		out = MustMarshalToBytes(NewServerRespMessages(s.IndexBefore, m))
		return
	case TypeClientFetchTitle:
		messageType = ws.MessageBinary
		out = MustMarshalToBytes(NewServerRespTitle(Rooms[roomIdx].Title))
		return
	case TypeClientSendMessage:
		var message ClientSendMessage
		// 因为之前已经unmarshal一次了, 所以这里可以不检查
		json.Unmarshal(data, &message)

		// 不允许用户那边发来的消息是空的, 这里做了判断, 如果不符合协议直接关闭
		if message.Message == "" {
			c.Close()
			return
		}

		roomMessage := &singleMessage{ID: MsgID, Time: time.Now(), Sender: username, Message: message.Message}
		Rooms[roomIdx].Messages = append(Rooms[roomIdx].Messages, roomMessage)
		MsgID++
		for _, v := range Conns {
			b, err := util.PackData(ws.MessageBinary, MustMarshalToBytes(NewServerRecvMessage(roomMessage)))
			if err != nil {
				panic(err)
			}
			v.Send(b)
		}
	case TypeHeartbeatMessage:
		messageType = ws.MessageBinary
		out = MustMarshalToBytes(NewHeartbeatMessage())
		return
	case TypeUnknown:
		log.Println("unknown packet, closing...")
		c.Close()
	}
	return
}

func (s *ChatServer) OnClose(c *gev.Connection) {
	log.Println("OnClose: ", c.PeerAddr())

	connID_, _ := c.Get(KeyConnID)
	connID := connID_.(int)
	username_, _ := c.Get(KeyUsername)
	username := username_.(string)
	authed_, _ := c.Get(KeyAuthed)
	authed := authed_.(bool)
	if authed {
		Users[username].Logined = false
	}
	delete(Conns, connID)

}
