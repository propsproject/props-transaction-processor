package state

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type State struct {
	context *processor.Context
}

var myLogger, _ = zap.NewProduction()
var logger = myLogger.Sugar()

func NewState(context *processor.Context) *State {

	logger = logger.With(
		zap.String("app", viper.GetString("app")),
		zap.String("name", viper.GetString("name")),
		zap.String("env", viper.GetString("environment")),
	)

	return &State{context: context}
}
