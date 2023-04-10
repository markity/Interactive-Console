package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	interactive "github.com/markity/Interactive-Console"

	"github.com/gorilla/websocket"
)

func main() {
	userInfo, err := ReadUsernamePassword("userinfo.json")
	if err != nil {
		fmt.Println("读取用户信息失败: " + err.Error())
		return
	}

	cfg := interactive.Config{
		Prompt:               '>',
		PromptStyle:          Style0,
		BlockInputAfterRun:   false,
		BlockInputAfterEnter: false,
		TraceAfterRun:        false,
		EventHandleMask:      interactive.EventMaskKeyUpWhenTrace | interactive.EventMaskTryToMoveUpper | interactive.EventMaskTryToMoveLower,
	}
	w := interactive.Run(cfg)

	ClearAndPrintHelp(w)

	// 0号协程用来跑事件循环
	cmdChan := w.GetCmdChan()
	eventChan := w.GetEventChan()
	var recvMsgChan chan interface{} = nil
	var writeMsgChan chan interface{} = nil
	var errChan chan string = nil

	// 保存页面的状态
	status := &Status{
		Connected:     false,
		RoomID:        -1,
		AtChatPage:    false,
		JustConnected: false,
		RoomTitle:     "",
		HeartbeatCnt:  0,
		Conn:          nil,
		Msgs:          nil,
	}

	// 事件循环
	for {
		select {
		// 解析命令
		case cmd_ := <-cmdChan:
			cmd, cmdLength := GetCmds(cmd_)
			switch {
			case cmdLength == 0:
				status.AtChatPage = false
				w.SetTrace(false)
				ClearAndPrintHelp(w)
				continue
			// list: 查询全部聊天室, 这是http请求
			case cmdLength == 1 && cmd[0] == "list":
				status.AtChatPage = false
				w.SetTrace(false)
				w.Clear()
				w.SendLineBackWithColor(Style2, "正在查询...")
				w.SetBlockInput(true)
				w.SetPrompt(nil, &Style3)
				listJSON, err := DoListRequest(userInfo)
				if err != nil {
					w.Clear()
					w.SendLineBackWithColor(Style1, "发送请求失败: "+err.Error())
					w.SetBlockInput(false)
					w.SetPrompt(nil, &Style0)
					continue
				}
				if listJSON.Code != 10000 {
					w.Clear()
					w.SendLineBackWithColor(Style1, "状态码"+fmt.Sprint(listJSON.Code)+": "+listJSON.Msg)
					w.SetBlockInput(false)
					w.SetPrompt(nil, &Style0)
					continue
				}
				w.Clear()
				w.SendLineBackWithColor(Style1, "所有的房间信息:")
				for _, v := range listJSON.RoomInfos {
					w.SendLineBackWithColor(Style2, fmt.Sprintf("%v: %v", v.Index, v.Title))
				}
				w.SetPrompt(nil, &Style0)
				w.SetBlockInput(false)
			case cmdLength == 2 && cmd[0] == "connect":
				status.AtChatPage = false
				w.SetTrace(false)
				roomNum, err := strconv.Atoi(cmd[1])
				if err != nil || roomNum <= 0 {
					w.Clear()
					w.SendLineBackWithColor(Style1, "不存在此房间号")
					continue
				}
				if status.Connected {
					w.Clear()
					w.SendLineBackWithColor(Style1, "你已经连接, 请先断开连接")
					continue
				}
				w.Clear()
				w.SetBlockInput(true)
				w.SetPrompt(nil, &Style3)
				w.SendLineBackWithColor(Style2, "正在连接到服务器...")
				conn, resp, err := DoDialChatRoom(userInfo, roomNum)
				if err != nil {
					if resp != nil {
						var b []byte
						b, err = io.ReadAll(resp.Body)
						resp.Body.Close()
						if err == nil {
							err = errors.New(string(b))
						}
					}
					w.Clear()
					w.SetPrompt(nil, &Style0)
					w.SendLineBackWithColor(Style1, "连接到服务器失败: "+err.Error())
					w.SetBlockInput(false)
					continue
				}
				w.Clear()
				w.SendLineFrontWithColor(Style2, "连接成功, 正在获取消息...")

				status = &Status{
					Connected:     true,
					RoomID:        roomNum,
					AtChatPage:    false,
					JustConnected: true,
					RoomTitle:     "",
					HeartbeatCnt:  0,
					Conn:          conn,
					Msgs:          nil,
				}

				recvMsgChan = make(chan interface{})
				writeMsgChan = make(chan interface{})
				errChan = make(chan string)
				go func() {
					go func() { writeMsgChan <- NewClientFetchTitle() }()
					go func() { writeMsgChan <- NewClientFetchMessages(0) }()
				}()

				// 创建reader和writer
				go func() {
					// reader
					go func() {
						for {
							_, binaryBytes, err := conn.ReadMessage()
							if err != nil {
								conn.Close()
								errChan <- err.Error()
								return
							}
							iface, messageType := ClientParseBytesToIface(binaryBytes)
							if messageType == TypeUnknown {
								conn.Close()
								errChan <- "协议错误, 无法解析对端的消息" + string(binaryBytes)
								return
							}
							go func() {
								recvMsgChan <- iface
							}()
						}
					}()

					// writer
					go func() {
						for v := range writeMsgChan {
							// error消息由reader发, writer不用发了, 不然管道关闭了就难受了
							err = conn.WriteMessage(websocket.BinaryMessage, MustMarshalToBytes(v))
							if err != nil {
								return
							}
						}
					}()
				}()
				continue

			// 查询连接状态
			case cmdLength == 1 && cmd[0] == "status":
				status.AtChatPage = false
				w.SetTrace(false)
				w.Clear()
				if !status.Connected {
					w.SendLineBackWithColor(Style2, "当前状态: 未连接")
				} else {
					w.SendLineBackWithColor(Style2, fmt.Sprintf("当前状态: 连接到聊天室%d(%s)", status.RoomID, status.RoomTitle))
				}
			case cmdLength == 1 && cmd[0] == "chatpage":
				if status.Connected {
					status.AtChatPage = true
					w.SetTrace(true)
					w.Clear()
					top, _ := IsAlreayTop(status.Msgs)
					if top {
						w.SendLineBackWithColor(Style1, "没有更多消息了...")
					} else {
						w.SendLineBackWithColor(Style1, "按上键显示更多...")
					}
					for _, v := range status.Msgs {
						w.SendLineBackWithColor(Style2, fmt.Sprintf("%v[%v]: %v", v.Sender, v.Time.Format("15:04:05, 2006"), v.Message))
					}
				} else {
					w.Clear()
					w.SendLineBackWithColor(Style1, "你还没连接到任何聊天室")
				}
			case cmdLength != 1 && cmd[0] == "send":
				msg := strings.TrimSpace(strings.TrimSpace(cmd_)[4:])
				w.Clear()
				if !status.Connected {
					w.SendLineFrontWithColor(Style1, "你还没有连接到任何聊天室")
				} else {
					status.AtChatPage = true
					w.SetTrace(true)
					top, _ := IsAlreayTop(status.Msgs)
					if top {
						w.SendLineBackWithColor(Style1, "没有更多消息了...")
					} else {
						w.SendLineBackWithColor(Style1, "按上键显示更多...")
					}
					for _, v := range status.Msgs {
						w.SendLineBackWithColor(Style2, fmt.Sprintf("%v[%v]: %v", v.Sender, v.Time.Format("15:04:05, 2006"), v.Message))
					}
					go func() { writeMsgChan <- NewClientSendMessage(msg) }()
				}
			case cmdLength == 1 && cmd[0] == "quit":
				goto out
			case cmdLength == 1 && cmd[0] == "disconnect":
				status.AtChatPage = false
				w.SetTrace(false)
				w.Clear()
				if status.Connected {
					w.SendLineBackWithColor(Style2, "已断开连接")
					status.Conn.Close()
					errChan = nil
					recvMsgChan = nil
					writeMsgChan = nil
					status = &Status{
						Connected:     false,
						RoomID:        -1,
						AtChatPage:    false,
						JustConnected: false,
						RoomTitle:     "",
						HeartbeatCnt:  0,
						Conn:          nil,
						Msgs:          nil,
					}
				} else {
					w.SendLineBackWithColor(Style1, "你还没有连接到任何聊天室")
				}
			default:
				ClearAndPrintHelp(w)
				continue
			}
		// 解析事件
		case ev := <-eventChan:
			if status.Connected && !status.JustConnected {
				switch ev.(type) {
				case *interactive.EventTypeUpWhenTrace:
					w.SetTrace(false)
				case *interactive.EventTryToGetLower:
					w.SetTrace(true)
				case *interactive.EventTryToGetUpper:
					alreadyTop, i := IsAlreayTop(status.Msgs)
					if alreadyTop {
						continue
					}
					go func() { writeMsgChan <- NewClientFetchMessages(i) }()
				}
			}
		// 消息通知
		case msg := <-recvMsgChan:
			switch gotMsg := msg.(type) {
			case *ServerRecvMessage:
				status.Msgs = append(status.Msgs, gotMsg.Message)
				if status.AtChatPage {
					w.SendLineBackWithColor(Style2, fmt.Sprintf("%v[%v]: %v", gotMsg.Message.Sender, gotMsg.Message.Time.Format("15:04:05, 2006"), gotMsg.Message.Message))
				}
			case *ServerRespMessages:
				if status.JustConnected {
					status.Msgs = gotMsg.Messages
					w.Clear()
					top, _ := IsAlreayTop(status.Msgs)
					if top {
						w.SendLineBackWithColor(Style1, "没有更多消息了...")
					} else {
						w.SendLineBackWithColor(Style1, "按上键显示更多...")
					}
					for _, v := range status.Msgs {
						w.SendLineBackWithColor(Style2, fmt.Sprintf("%v[%v]: %v", v.Sender, v.Time.Format("15:04:05, 2006"), v.Message))
					}
					if status.RoomTitle != "" {
						status.JustConnected = false
						status.AtChatPage = true
						w.SetBlockInput(false)
						w.SetPrompt(nil, &Style0)
						w.SetTrace(true)
					}
				} else {
					newMsgs := make([]*singleMessage, 0, 10+len(status.Msgs))
					newMsgs = append(newMsgs, gotMsg.Messages...)
					newMsgs = append(newMsgs, status.Msgs...)
					status.Msgs = newMsgs
					if status.AtChatPage {
						w.Clear()
						top, _ := IsAlreayTop(newMsgs)
						if top {
							w.SendLineBackWithColor(Style1, "没有更多消息了...")
						} else {
							w.SendLineBackWithColor(Style1, "按上键显示更多...")
						}
						for _, v := range newMsgs {
							w.SendLineBackWithColor(Style2, fmt.Sprintf("%v[%v]: %v", v.Sender, v.Time.Format("15:04:05, 2006"), v.Message))
						}
					}
				}
			case *ServerRespTitle:
				status.RoomTitle = gotMsg.Title
				if status.Msgs != nil {
					status.JustConnected = false
					status.AtChatPage = true
					w.SetBlockInput(false)
					w.SetPrompt(nil, &Style0)
					w.SetTrace(true)
				}
			case *HeartbeatMessage:
				status.HeartbeatCnt = 0
			}
		// 错误通知
		case v := <-errChan:
			recvMsgChan = nil
			writeMsgChan = nil
			errChan = nil
			w.Clear()
			w.SendLineFrontWithColor(Style1, "连接出现错误, 断开连接: "+v)
			w.SetPrompt(nil, &Style0)
			w.SetBlockInput(false)
			status = &Status{
				Connected:     false,
				RoomID:        -1,
				AtChatPage:    false,
				JustConnected: false,
				RoomTitle:     "",
				HeartbeatCnt:  0,
				Conn:          nil,
				Msgs:          nil,
			}
		// 发送心跳包, 进行心跳检测
		case <-time.Tick(time.Second * 3):
			if !status.Connected {
				continue
			} else {
				go func() {
					writeMsgChan <- NewHeartbeatMessage()
				}()
				if status.HeartbeatCnt == 3 {
					recvMsgChan = nil
					errChan = nil
					w.Clear()
					w.SendLineFrontWithColor(Style1, "连接心跳超时, 自动断开连接")
					w.SetBlockInput(false)
					status.Conn.Close()
					status = &Status{
						Connected:     false,
						RoomID:        -1,
						AtChatPage:    false,
						JustConnected: false,
						RoomTitle:     "",
						HeartbeatCnt:  0,
						Conn:          nil,
						Msgs:          nil,
					}
				}
				status.HeartbeatCnt++
			}
		}
	}
out:
	w.Stop()
}
