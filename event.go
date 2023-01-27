package interactive

import (
	"time"
)

type sendEvent struct {
	when time.Time
	data []rune
}

func (me *sendEvent) When() time.Time {
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

// TODO 是否支持左右移动?
// type gotoLeftEvent struct {
// 	when time.Time
// }

// func (me *gotoLeftEvent) When() time.Time {
// 	return me.when
// }

// type gotoRightEvent struct {
// 	when time.Time
// }

// func (me *gotoRightEvent) When() time.Time {
// 	return me.when
// }
