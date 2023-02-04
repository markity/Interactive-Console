package interactive

import (
	"time"
)

// 导出事件

// 用户按上键
const EventMaskKeyUp = 1

// 按下键
const EventMaskKeyDown = 2 << 0

// 已经在顶端, 此时用户按上键
const EventMaskTryToMoveUpper = 2 << 1

// 已经在底端, 此时用户按下键
const EventMaskTryToMoveLower = 2 << 2

// 用户在trace模式下按上下键
const EventMaskKeyUpWhenTrace = 2 << 3
const EventMaskKeyDownWhenTrace = 2 << 4

// 一些控制案按键, 为什么不提供Ctrl+M? 因为它就是Enter, 这个键位有特殊用途, 因此不提供使用
const EventMaskKeyCtrlSpace = 2 << 5
const EventMaskKeyCtrlA = 2 << 6
const EventMaskKeyCtrlB = 2 << 7
const EventMaskKeyCtrlC = 2 << 8
const EventMaskKeyCtrlD = 2 << 9
const EventMaskKeyCtrlE = 2 << 10
const EventMaskKeyCtrlF = 2 << 11
const EventMaskKeyCtrlG = 2 << 12
const EvnetMaskKeyCtrlH = 2 << 13
const EventMaskKeyCtrlI = 2 << 14
const EventMaskKeyCtrlJ = 2 << 15
const EventMaskKeyCtrlK = 2 << 16
const EventMaskKeyCtrlL = 2 << 17
const EventMaskKeyCtrlN = 2 << 18
const EventMaskKeyCtrlO = 2 << 19
const EventMaskKeyCtrlP = 2 << 20
const EventMaskKeyCtrlQ = 2 << 21
const EventMaskKeyCtrlR = 2 << 22
const EventMaskKeyCtrlS = 2 << 23
const EventMaskKeyCtrlT = 2 << 24
const EventMaskKeyCtrlU = 2 << 25
const EventMaskKeyCtrlV = 2 << 26
const EventMaskKeyCtrlW = 2 << 27
const EventMaskKeyCtrlX = 2 << 28
const EventMaskKeyCtrlY = 2 << 29
const EventMaskKeyCtrlZ = 2 << 30

const EventMaskWindowResize = 2 << 31

// 上移事件
type EventMoveUp struct {
	When                 time.Time
	LineOffsetBeforeMove int
}

// 下移事件
type EventMoveDown struct {
	When                 time.Time
	LineOffsetBeforeMove int
}

// 如果已经在最顶端, 还按上键, 那么产生这个事件
// 可以用来做聊天软件的查看历史消息功能
type EventTryToGetUpper struct {
	When time.Time
}

// 如果已经在最底端, 还按下键, 那么产生这个事件
// 可以指示程序输出更多
type EventTryToGetLower struct {
	When time.Time
}

// 在trace状态时按上键
type EventTypeUpWhenTrace struct {
	When time.Time
}

// 在trace状态时按下键
type EventTypeDownWhenTrace struct {
	When time.Time
}

type EventKeyCtrlSpace struct {
	When time.Time
}

type EventKeyCtrlA struct {
	When time.Time
}

type EventKeyCtrlB struct {
	When time.Time
}

type EventKeyCtrlC struct {
	When time.Time
}

type EventKeyCtrlD struct {
	When time.Time
}

type EventKeyCtrlE struct {
	When time.Time
}

type EventKeyCtrlF struct {
	When time.Time
}

type EventKeyCtrlG struct {
	When time.Time
}

type EventKeyCtrlH struct {
	When time.Time
}

type EventKeyCtrlI struct {
	When time.Time
}

type EventKeyCtrlJ struct {
	When time.Time
}

type EventKeyCtrlK struct {
	When time.Time
}

type EventKeyCtrlL struct {
	When time.Time
}

type EventKeyCtrlN struct {
	When time.Time
}

type EventKeyCtrlO struct {
	When time.Time
}

type EventKeyCtrlP struct {
	When time.Time
}

type EventKeyCtrlQ struct {
	When time.Time
}

type EventKeyCtrlR struct {
	When time.Time
}

type EventKeyCtrlS struct {
	When time.Time
}

type EventKeyCtrlT struct {
	When time.Time
}

type EventKeyCtrlU struct {
	When time.Time
}

type EventKeyCtrlV struct {
	When time.Time
}

type EventKeyCtrlW struct {
	When time.Time
}

type EventKeyCtrlX struct {
	When time.Time
}

type EventKeyCtrlY struct {
	When time.Time
}

type EventKeyCtrlZ struct {
	When time.Time
}

type EventWindowResize struct {
	Height int
	Width  int
	When   time.Time
}

// 内部事件

type stopEvent struct {
	when time.Time
}

func (ev *stopEvent) When() time.Time {
	return ev.when
}

type clearEvent struct {
	when time.Time
}

func (me *clearEvent) When() time.Time {
	return me.when
}

type gotoBottomEvent struct {
	when time.Time
}

func (me *gotoBottomEvent) When() time.Time {
	return me.when
}

type gotoTopEvent struct {
	when time.Time
}

func (me *gotoTopEvent) When() time.Time {
	return me.when
}

type setBlockInputAfterEnterEvent struct {
	when time.Time
	data bool
}

func (me *setBlockInputAfterEnterEvent) When() time.Time {
	return me.when
}

type setTraceEvent struct {
	when time.Time
	data bool
}

func (me *setTraceEvent) When() time.Time {
	return me.when
}

type setBlockInputChangeEvent struct {
	when time.Time
	data bool
}

func (me *setBlockInputChangeEvent) When() time.Time {
	return me.when
}

type gotoLeftEvent struct {
	when time.Time
}

func (me *gotoLeftEvent) When() time.Time {
	return me.when
}

type gotoLineEvent struct {
	when time.Time
	data int
}

func (me *gotoLineEvent) When() time.Time {
	return me.when
}

type gotoNextLineEvent struct {
	when time.Time
}

func (me *gotoNextLineEvent) When() time.Time {
	return me.when
}

type gotoPreviousLineEvent struct {
	when time.Time
}

func (me *gotoPreviousLineEvent) When() time.Time {
	return me.when
}

type popFrontLineEvent struct {
	when time.Time
}

func (me *popFrontLineEvent) When() time.Time {
	return me.when
}

type popBackLineEvent struct {
	when time.Time
}

func (me *popBackLineEvent) When() time.Time {
	return me.when
}

type sendLineBackWithColorEvent struct {
	when time.Time
	data []interface{}
}

func (me *sendLineBackWithColorEvent) When() time.Time {
	return me.when
}

type sendLineFrontWithColorEvent struct {
	when time.Time
	data []interface{}
}

func (me *sendLineFrontWithColorEvent) When() time.Time {
	return me.when
}

type setPromptEvent struct {
	when      time.Time
	dataRune  *rune
	dataStyle *StyleAttr
}

func (me *setPromptEvent) When() time.Time {
	return me.when
}

type getWindowSizeEventResp struct {
	height int
	width  int
}

type getWindowSizeEvent struct {
	when time.Time
	resp chan *getWindowSizeEventResp
}

func (me *getWindowSizeEvent) When() time.Time {
	return me.when
}

// type gotoRightEvent struct {
// 	when time.Time
// }

// func (me *gotoRightEvent) When() time.Time {
// 	return me.when
// }
