package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	interactive "github.com/markity/Interactive-Console"

	"github.com/gorilla/websocket"
)

var ListURL = "http://localhost:8000/list"
var ChatURL = "ws://localhost:8001/chat/"

// 粗体, 粉红色, 代表prompt能输入
var Style0 interactive.StyleAttr

// 粗体, 红色, 代表有错误
var Style1 interactive.StyleAttr

// 非粗体, 蓝色, 代表正常信息
var Style2 interactive.StyleAttr

// 非粗体, 分红色, 代表prompt不能输入
var Style3 interactive.StyleAttr

func init() {
	Style0 = interactive.GetDefaultSytleAttr()
	Style0.Bold = true
	Style0.Underline = true
	Style0.Foreground = interactive.Color162

	Style1 = interactive.GetDefaultSytleAttr()
	Style1.Foreground = interactive.Color202
	Style1.Bold = true

	Style2 = interactive.GetDefaultSytleAttr()
	Style2.Foreground = interactive.Color81

	Style3 = interactive.GetDefaultSytleAttr()
	Style3.Bold = false
	Style3.Underline = true
	Style3.Dim = true
	Style3.Foreground = interactive.Color162
}

type AuthInfo struct {
	Username string
	Password string
}

func ReadUsernamePassword(path string) (*AuthInfo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	uname, ok := m["username"].(string)
	if !ok {
		return nil, errors.New("lack of username field")
	}

	pword, ok := m["password"].(string)
	if !ok {
		return nil, errors.New("lack of password field")
	}

	return &AuthInfo{Username: uname, Password: pword}, nil
}

func ClearAndPrintHelp(w *interactive.Win) {
	w.Clear()
	w.SendLineBackWithColor(Style1, "你现在在主界面, 全部指令:")
	w.SendLineBackWithColor(Style2, "                list: 查询所有聊天室及在线人数")
	w.SendLineBackWithColor(Style2, "   connect <room id>: 连接到服务器, 进行聊天")
	w.SendLineBackWithColor(Style2, "              status: 查看当前连接状态")
	w.SendLineBackWithColor(Style2, "            chatpage: 切换到聊天界面")
	w.SendLineBackWithColor(Style2, "      send <content>: 发送信息")
	w.SendLineBackWithColor(Style2, "          disconnect: 关闭连接, 返回主界面")
	w.SendLineBackWithColor(Style2, "                quit: 退出程序")
}

func DoListRequest(userInfo *AuthInfo) (*ListRequestResp, error) {
	l := ListRequestResp{}
	req, err := http.NewRequest("GET", ListURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("username", userInfo.Username)
	req.Header.Set("password", userInfo.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func DoDialChatRoom(userInfo *AuthInfo, roomNum int) (*websocket.Conn, *http.Response, error) {
	hdr := http.Header{}
	hdr.Set("username", userInfo.Username)
	hdr.Set("password", userInfo.Password)
	websocket.DefaultDialer.HandshakeTimeout = time.Second * 5
	return websocket.DefaultDialer.Dial(ChatURL+fmt.Sprint(roomNum), hdr)
}

// 保存界面的状态
type Status struct {
	Connected     bool
	RoomID        int
	AtChatPage    bool
	JustConnected bool
	RoomTitle     string
	HeartbeatCnt  int
	Conn          *websocket.Conn
	Msgs          []*singleMessage
}

func GetCmds(str string) ([]string, int) {
	cmds := strings.Fields(strings.TrimSpace(str))
	cmdLength := len(cmds)
	return cmds, cmdLength
}

func IsAlreayTop(msg []*singleMessage) (bool, int) {
	if len(msg) == 0 {
		return true, -1
	}
	if msg[0].ID == 1 {
		return true, -1
	}
	return false, msg[0].ID
}
