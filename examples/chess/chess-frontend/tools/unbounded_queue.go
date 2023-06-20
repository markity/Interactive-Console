package tools

import "sync"

type unboundedQueen struct {
	queen []interface{}
	l     sync.Mutex
	cond  *sync.Cond
}

func NewUnboundedQueen() *unboundedQueen {
	uq := unboundedQueen{}
	uq.queen = make([]interface{}, 0)
	uq.cond = sync.NewCond(&uq.l)
	return &uq
}

func (uq *unboundedQueen) Push(i interface{}) {
	uq.l.Lock()
	uq.queen = append(uq.queen, i)
	uq.l.Unlock()
	uq.cond.Broadcast()
}

func (uq *unboundedQueen) PopBlock() interface{} {
	uq.l.Lock()
	for {
		if len(uq.queen) == 0 {
			uq.cond.Wait()
		} else {
			ret := uq.queen[0]
			uq.queen = uq.queen[1:]
			uq.l.Unlock()
			return ret
		}
	}
}
