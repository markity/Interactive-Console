package interactive

import (
	"time"
)

// 导出事件

const EventMaskMoveUp = 1
const EventMaskMoveDown = 2
const EventMaskTryToGetUpper = 4
const EventMaskTryToGetLower = 8
const EventMaskTypeUpWhenTrace = 16
const EventMaskTypeDownWhenTrace = 16

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

// 内部事件
type sendEvent struct {
	when time.Time
	data []rune
}

func (me *sendEvent) When() time.Time {
	return me.when
}

type sendFrontEvent struct {
	when time.Time
	data []rune
}

func (me *sendFrontEvent) When() time.Time {
	return me.when
}

type endEvent struct {
	when time.Time
}

func (ev *endEvent) When() time.Time {
	return ev.when
}

type clearEvent struct {
	when time.Time
}

func (me *clearEvent) When() time.Time {
	return me.when
}

type gotoButtomEvent struct {
	when time.Time
}

func (me *gotoButtomEvent) When() time.Time {
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

// type gotoRightEvent struct {
// 	when time.Time
// }

// func (me *gotoRightEvent) When() time.Time {
// 	return me.when
// }
