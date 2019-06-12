package rpc

import (
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
)

type IteratorCB func(earning *pending_props_pb.Earning) error

type EarningsIterator struct {
	Earnings Earnings
	current  int
	err      error
}

func (e *EarningsIterator) Next(cb IteratorCB) (bool) {
	if e.err != nil || e.current == len(e.Earnings){
		return false
	}

	e.err  = cb(&e.Earnings[e.current])
	e.current++

	if e.current < len(e.Earnings) {
		return true
	}

	return false
}

type Earnings []pending_props_pb.Earning

func (e *EarningsIterator) CurrentValue() Earnings {
	return e.Earnings
}

func (e *EarningsIterator) Err() error {
	return e.err
}

func (e Earnings) NewIterator() *EarningsIterator {
	return &EarningsIterator{e, 0, nil}
}

func (e Earnings) Len() int {
	return len(e)
}

func (e Earnings) Less(i, j int) bool {
	return e[i].GetDetails().GetTimestamp() < e[j].GetDetails().GetTimestamp()
}

func (e Earnings) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
