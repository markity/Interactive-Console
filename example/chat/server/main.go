package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

// connection lifecycle
// OnConnect() -> OnRequest() -> OnHeader() -> OnMessage() -> OnClose()

// Users和Rooms写死了, 不用加锁
var Users = map[string]*UserData{}
var Rooms = map[int]*RoomData{}
var MsgID int = 1
var ConnID int = 1
var Conns map[int]*gev.Connection

func main() {
	// 开gev websocket服务
	go func() {
		handler := &ChatServer{}

		wsUpgrader := &ws.Upgrader{}
		wsUpgrader.OnRequest = func(c *gev.Connection, uri []byte) error {
			log.Println("OnRequest: ", string(uri))
			for _, v := range Rooms {
				if string(uri) == fmt.Sprintf("/chat/%d", v.ID) {
					c.Set(KeyUri, string(uri))
					c.Set(KeyRoomIndex, v.ID)
					c.Set(KeyConnID, ConnID)
					Conns[ConnID] = c
					ConnID++
					return nil
				}
			}
			return errors.New("没有这个房间")
		}

		wsUpgrader.OnBeforeUpgrade = func(c *gev.Connection) (header ws.HandshakeHeader, err error) {
			headerIface, _ := c.Get(KeyRequestHeader)
			httpHeader := headerIface.(http.Header)
			u, exists := Users[httpHeader.Get("username")]
			if !exists || u.Password != httpHeader.Get("password") {
				c.Close()
				return nil, errors.New("鉴权失败")
			}
			c.Set(KeyUsername, httpHeader.Get("username"))
			if u.Logined {
				c.Close()
				c.Set(KeyAuthed, false)
				return nil, errors.New("用户已经登陆, 请勿重复登陆")
			}
			c.Set(KeyAuthed, true)
			u.Logined = true
			return nil, nil
		}

		wsUpgrader.OnHeader = func(c *gev.Connection, key, value []byte) error {
			log.Println("OnHeader: ", string(key), string(value))

			var header http.Header
			_header, ok := c.Get(KeyRequestHeader)
			if ok {
				header = _header.(http.Header)
			} else {
				header = make(http.Header)
			}
			header.Set(string(key), string(value))

			c.Set(KeyRequestHeader, header)
			return nil
		}

		s, err := NewWebSocketServer(handler, wsUpgrader,
			gev.Network("tcp"),
			gev.Address("localhost:8001"),
			gev.NumLoops(1))
		if err != nil {
			panic(err)
		}

		s.Start()
	}()

	go GoHttpService()
	select {}
}
