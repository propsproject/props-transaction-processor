package state

import (
	"github.com/propsproject/sawtooth-go-sdk/logging"
	"github.com/propsproject/sawtooth-go-sdk/processor"
)

type State struct {
	context *processor.Context
}

var logger = logging.Get()

func NewState(context *processor.Context) *State {
	return &State{context: context}
}
