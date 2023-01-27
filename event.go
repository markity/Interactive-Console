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

type blockInputChangeEvent struct {
	when time.Time
	data bool
}

func (me *blockInputChangeEvent) When() time.Time {
	return me.when
}

type clearEvent struct {
	when time.Time
}

func (me *clearEvent) When() time.Time {
	return me.when
}
