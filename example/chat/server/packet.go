package main

import (
	"encoding/json"
	"time"
)

// 封装了客户端/服务端的消息包相关结构体及功能函数

type MessageType int

const TypeClientFetchTitle MessageType = 1
const TypeClientFetchMessages MessageType = 2
const TypeClientSendMessage MessageType = 3
const TypeServerRespTitle MessageType = 4
const TypeServerRespMessages MessageType = 5
const TypeServerRecvMessage MessageType = 6
const TypeHeartbeatMessage MessageType = 7
const TypeUnknown MessageType = 8

type PacketHeader struct {
	Type string `json:"type"`
}

// 1
type ClientFetchTitle struct {
	PacketHeader
}

func NewClientFetchTitle() *ClientFetchTitle {
	return &ClientFetchTitle{
		PacketHeader: PacketHeader{Type: "client_fetch_title"},
	}
}

// 2
type ClientFetchMessages struct {
	PacketHeader
	IndexBefore int `json:"index_before"`
}

func NewClientFetchMessages(idx_before int) *ClientFetchMessages {
	return &ClientFetchMessages{PacketHeader: PacketHeader{Type: "client_fetch_messages"}, IndexBefore: idx_before}
}

// 3
type ClientSendMessage struct {
	PacketHeader
	Message string `json:"message"`
}

func NewClientSendMessage(msg string) *ClientSendMessage {
	return &ClientSendMessage{PacketHeader: PacketHeader{Type: "client_send_message"}, Message: msg}
}

// 4
type ServerRespTitle struct {
	PacketHeader
	Title string `json:"title"`
}

func NewServerRespTitle(title string) *ServerRespTitle {
	return &ServerRespTitle{PacketHeader: PacketHeader{Type: "server_resp_title"}, Title: title}
}

type singleMessage struct {
	ID      int       `json:"id"`
	Sender  string    `json:"sender"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

// 5
type ServerRespMessages struct {
	PacketHeader
	IndexBefore int              `json:"index_before"`
	Messages    []*singleMessage `json:"messages"`
}

func NewServerRespMessages(idx_before int, messages []*singleMessage) *ServerRespMessages {
	return &ServerRespMessages{PacketHeader: PacketHeader{Type: "server_resp_messages"}, IndexBefore: idx_before, Messages: messages}
}

// 6
type ServerRecvMessage struct {
	PacketHeader
	Message *singleMessage `json:"message"`
}

func NewServerRecvMessage(msg *singleMessage) *ServerRecvMessage {
	return &ServerRecvMessage{PacketHeader: PacketHeader{Type: "server_recv_message"}, Message: msg}
}

// 7
type HeartbeatMessage struct {
	PacketHeader
}

func NewHeartbeatMessage() *HeartbeatMessage {
	return &HeartbeatMessage{PacketHeader: PacketHeader{Type: "heartbeat_message"}}
}

// 工具函数, 将interface{}转化成json bytes, 如果失败直接panic, 需要保证传入的结构体合法
func MustMarshalToBytes(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return b
}

// 工具函数, 侦查消息类型
func DetectMessageType(b []byte) MessageType {
	var header PacketHeader
	if err := json.Unmarshal(b, &header); err != nil {
		return TypeUnknown
	}

	switch header.Type {
	case "client_fetch_title":
		return TypeClientFetchTitle
	case "client_fetch_messages":
		return TypeClientFetchMessages
	case "client_send_message":
		return TypeClientSendMessage
	case "server_resp_title":
		return TypeServerRespTitle
	case "server_recv_message":
		return TypeServerRecvMessage
	case "heartbeat_message":
		return TypeHeartbeatMessage
	default:
		return TypeUnknown
	}
}
