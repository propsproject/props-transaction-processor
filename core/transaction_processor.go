package core

import (
	"github.com/propsproject/sawtooth-go-sdk/processor"
	"syscall"
)

func NewTransactionProcessor(validator string) *processor.TransactionProcessor {
	earningsHandler := NewTransactionHandler()

	processor := processor.NewTransactionProcessor(validator)
	processor.AddHandler(earningsHandler)
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)

	return processor
}